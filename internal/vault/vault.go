package vault

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	va "github.com/hashicorp/vault-client-go"
	"github.com/hashicorp/vault-client-go/schema"
)

var (
	ErrItemAlreadyExist = errors.New("item already exist")
	ErrKeyNotFound      = errors.New("key not found")
)

type Vault[T any] interface {
	Add(ctx context.Context, name string, model T) error
	Get(ctx context.Context, name string) (*T, error)
	List(ctx context.Context) (map[string]T, error)
	Remove(ctx context.Context, name string) error
	Update(ctx context.Context, name string, model T) error
	Ping(ctx context.Context) error
}

type Client struct {
	*va.Client
	engine string
}

type vault[T any] struct {
	client *Client
	key    string
}

func New(ctx context.Context, addr, token, engine, description string) (*Client, error) {
	c, err := va.New(va.WithAddress(addr))
	if err != nil {
		return nil, err
	}

	if err = c.SetToken(token); err != nil {
		return nil, err
	}

	if _, err = c.System.MountsEnableSecretsEngine(ctx, engine, schema.MountsEnableSecretsEngineRequest{
		Description: description,
		Type:        "kv",
	}); err != nil && !strings.Contains(err.Error(), "path is already in use") {
		return nil, err
	}

	return &Client{Client: c, engine: engine}, nil
}

func NewRepository[T any](client *Client, key string) Vault[T] {
	return &vault[T]{
		client: client,
		key:    key,
	}
}

func (v *vault[T]) Add(ctx context.Context, name string, model T) error {
	items, err := v.List(ctx)
	if err != nil && !errors.Is(err, ErrKeyNotFound) {
		return err
	}

	if items == nil {
		items = make(map[string]T)
	}

	if _, ok := items[name]; ok {
		return ErrItemAlreadyExist
	}

	items[name] = model

	if _, err = v.client.Secrets.KvV2Write(ctx, v.key, schema.KvV2WriteRequest{
		Data: map[string]any{
			"items": items,
		},
	}, va.WithMountPath(v.client.engine)); err != nil {
		return err
	}

	return nil
}

func (v *vault[T]) Get(ctx context.Context, name string) (*T, error) {
	items, err := v.List(ctx)
	if err != nil {
		return nil, err
	}

	for k, i := range items {
		if k == name {
			return &i, nil
		}
	}

	return nil, ErrKeyNotFound
}

func (v *vault[T]) List(ctx context.Context) (map[string]T, error) {
	items := make(map[string]T)

	res, err := v.client.Secrets.KvV2Read(ctx, v.key, va.WithMountPath(v.client.engine))
	if err != nil && err.(*va.ResponseError).StatusCode != http.StatusNotFound {
		return nil, err
	}

	if err != nil && err.(*va.ResponseError).StatusCode == http.StatusNotFound {
		return nil, ErrKeyNotFound
	}

	if len(res.Data.Data) == 0 || len(res.Data.Data["items"].(map[string]any)) == 0 {
		return nil, ErrKeyNotFound
	}

	m, err := json.Marshal(res.Data.Data["items"])
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(m, &items); err != nil {
		return nil, err
	}

	return items, nil
}

func (v *vault[T]) Remove(ctx context.Context, name string) error {
	items, err := v.List(ctx)
	if err != nil {
		return err
	}

	if _, ok := items[name]; !ok {
		return ErrKeyNotFound
	}

	delete(items, name)

	if _, err = v.client.Secrets.KvV2Write(ctx, v.key, schema.KvV2WriteRequest{
		Data: map[string]any{
			"items": items,
		},
	}, va.WithMountPath(v.client.engine)); err != nil {
		return err
	}

	return nil
}

func (v *vault[T]) Update(ctx context.Context, name string, model T) error {
	items, err := v.List(ctx)
	if err != nil {
		return err
	}

	if _, ok := items[name]; !ok {
		return ErrKeyNotFound
	}

	items[name] = model

	if _, err = v.client.Secrets.KvV2Write(ctx, v.key, schema.KvV2WriteRequest{
		Data: map[string]any{
			"items": items,
		},
	}, va.WithMountPath(v.client.engine)); err != nil {
		return err
	}

	return nil
}

func (v *vault[T]) Ping(ctx context.Context) error {
	_, err := v.client.System.ReadHealthStatus(ctx, va.WithMountPath(v.client.engine))
	if err != nil {
		return err
	}

	return nil
}
