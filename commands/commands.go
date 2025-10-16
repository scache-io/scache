package commands

import (
	"errors"
	"fmt"
	"time"

	"github.com/scache-io/scache/interfaces"
	"github.com/scache-io/scache/types"
)

// BaseCommand 基础命令实现
type BaseCommand struct {
	name string
}

func (c *BaseCommand) Name() string {
	return c.name
}

func (c *BaseCommand) Validate(args []interface{}) error {
	// 基础验证，子类可以覆盖
	return nil
}

// SetCommand SET 命令
type SetCommand struct {
	BaseCommand
}

func NewSetCommand() *SetCommand {
	return &SetCommand{BaseCommand: BaseCommand{name: "SET"}}
}

func (c *SetCommand) Execute(ctx *interfaces.Context) error {
	if len(ctx.Args) < 2 {
		return errors.New("SET requires at least 2 arguments")
	}

	key, ok := ctx.Args[0].(string)
	if !ok {
		return errors.New("key must be string")
	}

	value := ctx.Args[1]

	var ttl time.Duration
	if len(ctx.Args) >= 3 {
		if t, ok := ctx.Args[2].(time.Duration); ok {
			ttl = t
		}
	}

	// 创建字符串对象
	obj := types.NewStringObject(fmt.Sprintf("%v", value), ttl)
	return ctx.Storage.Set(key, obj)
}

// GetCommand GET 命令
type GetCommand struct {
	BaseCommand
}

func NewGetCommand() *GetCommand {
	return &GetCommand{BaseCommand: BaseCommand{name: "GET"}}
}

func (c *GetCommand) Execute(ctx *interfaces.Context) error {
	if len(ctx.Args) < 1 {
		return errors.New("GET requires at least 1 argument")
	}

	key, ok := ctx.Args[0].(string)
	if !ok {
		return errors.New("key must be string")
	}

	obj, found := ctx.Storage.Get(key)
	if !found {
		ctx.Result = nil
		return nil
	}

	if strObj, ok := obj.(*types.StringObject); ok {
		ctx.Result = strObj.Value()
	} else {
		return errors.New("key does not contain a string")
	}

	return nil
}

// DeleteCommand DEL 命令
type DeleteCommand struct {
	BaseCommand
}

func NewDeleteCommand() *DeleteCommand {
	return &DeleteCommand{BaseCommand: BaseCommand{name: "DEL"}}
}

func (c *DeleteCommand) Execute(ctx *interfaces.Context) error {
	if len(ctx.Args) < 1 {
		return errors.New("DEL requires at least 1 argument")
	}

	key, ok := ctx.Args[0].(string)
	if !ok {
		return errors.New("key must be string")
	}

	deleted := ctx.Storage.Delete(key)
	ctx.Result = deleted
	return nil
}

// ExistsCommand EXISTS 命令
type ExistsCommand struct {
	BaseCommand
}

func NewExistsCommand() *ExistsCommand {
	return &ExistsCommand{BaseCommand: BaseCommand{name: "EXISTS"}}
}

func (c *ExistsCommand) Execute(ctx *interfaces.Context) error {
	if len(ctx.Args) < 1 {
		return errors.New("EXISTS requires at least 1 argument")
	}

	key, ok := ctx.Args[0].(string)
	if !ok {
		return errors.New("key must be string")
	}

	exists := ctx.Storage.Exists(key)
	ctx.Result = exists
	return nil
}

// LPushCommand LPUSH 命令
type LPushCommand struct {
	BaseCommand
}

func NewLPushCommand() *LPushCommand {
	return &LPushCommand{BaseCommand: BaseCommand{name: "LPUSH"}}
}

func (c *LPushCommand) Execute(ctx *interfaces.Context) error {
	if len(ctx.Args) < 2 {
		return errors.New("LPUSH requires at least 2 arguments")
	}

	key, ok := ctx.Args[0].(string)
	if !ok {
		return errors.New("key must be string")
	}

	value := ctx.Args[1]

	var ttl time.Duration
	if len(ctx.Args) >= 3 {
		if t, ok := ctx.Args[2].(time.Duration); ok {
			ttl = t
		}
	}

	// 获取或创建列表对象
	obj, exists := ctx.Storage.Get(key)
	var listObj *types.ListObject

	if exists {
		if existingList, ok := obj.(*types.ListObject); ok {
			listObj = existingList
		} else {
			return errors.New("key exists but is not a list")
		}
	} else {
		listObj = types.NewListObject([]interface{}{value}, ttl)
		if err := ctx.Storage.Set(key, listObj); err != nil {
			return err
		}
		ctx.Result = 1
		return nil
	}

	// 添加元素到列表开头
	values := listObj.Values()
	newValues := make([]interface{}, len(values)+1)
	newValues[0] = value
	copy(newValues[1:], values)

	// 创建新的列表对象
	newListObj := types.NewListObject(newValues, ttl)
	if err := ctx.Storage.Set(key, newListObj); err != nil {
		return err
	}

	ctx.Result = len(newValues)
	return nil
}

// RPopCommand RPOP 命令
type RPopCommand struct {
	BaseCommand
}

func NewRPopCommand() *RPopCommand {
	return &RPopCommand{BaseCommand: BaseCommand{name: "RPOP"}}
}

func (c *RPopCommand) Execute(ctx *interfaces.Context) error {
	if len(ctx.Args) < 1 {
		return errors.New("RPOP requires at least 1 argument")
	}

	key, ok := ctx.Args[0].(string)
	if !ok {
		return errors.New("key must be string")
	}

	obj, exists := ctx.Storage.Get(key)
	if !exists {
		ctx.Result = nil
		return nil
	}

	listObj, ok := obj.(*types.ListObject)
	if !ok {
		return errors.New("key does not contain a list")
	}

	value, hasValue := listObj.Pop()
	if !hasValue {
		ctx.Result = nil
		return nil
	}

	// 更新存储的列表对象
	if err := ctx.Storage.Set(key, listObj); err != nil {
		return err
	}

	ctx.Result = value
	return nil
}

// HSetCommand HSET 命令
type HSetCommand struct {
	BaseCommand
}

func NewHSetCommand() *HSetCommand {
	return &HSetCommand{BaseCommand: BaseCommand{name: "HSET"}}
}

func (c *HSetCommand) Execute(ctx *interfaces.Context) error {
	if len(ctx.Args) < 3 {
		return errors.New("HSET requires at least 3 arguments")
	}

	key, ok := ctx.Args[0].(string)
	if !ok {
		return errors.New("key must be string")
	}

	field, ok := ctx.Args[1].(string)
	if !ok {
		return errors.New("field must be string")
	}

	value := ctx.Args[2]

	var ttl time.Duration
	if len(ctx.Args) >= 4 {
		if t, ok := ctx.Args[3].(time.Duration); ok {
			ttl = t
		}
	}

	// 获取或创建哈希对象
	obj, exists := ctx.Storage.Get(key)
	var hashObj *types.HashObject

	if exists {
		if existingHash, ok := obj.(*types.HashObject); ok {
			hashObj = existingHash
		} else {
			return errors.New("key exists but is not a hash")
		}
	} else {
		fields := make(map[string]interface{})
		fields[field] = value
		hashObj = types.NewHashObject(fields, ttl)
		if err := ctx.Storage.Set(key, hashObj); err != nil {
			return err
		}
		ctx.Result = 1
		return nil
	}

	// 设置字段
	hashObj.Set(field, value)

	// 重新设置以更新存储
	if err := ctx.Storage.Set(key, hashObj); err != nil {
		return err
	}

	ctx.Result = 1
	return nil
}

// HGetCommand HGET 命令
type HGetCommand struct {
	BaseCommand
}

func NewHGetCommand() *HGetCommand {
	return &HGetCommand{BaseCommand: BaseCommand{name: "HGET"}}
}

func (c *HGetCommand) Execute(ctx *interfaces.Context) error {
	if len(ctx.Args) < 2 {
		return errors.New("HGET requires at least 2 arguments")
	}

	key, ok := ctx.Args[0].(string)
	if !ok {
		return errors.New("key must be string")
	}

	field, ok := ctx.Args[1].(string)
	if !ok {
		return errors.New("field must be string")
	}

	obj, exists := ctx.Storage.Get(key)
	if !exists {
		ctx.Result = nil
		return nil
	}

	hashObj, ok := obj.(*types.HashObject)
	if !ok {
		return errors.New("key does not contain a hash")
	}

	value, found := hashObj.Get(field)
	if !found {
		ctx.Result = nil
	} else {
		ctx.Result = value
	}

	return nil
}

// TypeCommand TYPE 命令
type TypeCommand struct {
	BaseCommand
}

func NewTypeCommand() *TypeCommand {
	return &TypeCommand{BaseCommand: BaseCommand{name: "TYPE"}}
}

func (c *TypeCommand) Execute(ctx *interfaces.Context) error {
	if len(ctx.Args) < 1 {
		return errors.New("TYPE requires at least 1 argument")
	}

	key, ok := ctx.Args[0].(string)
	if !ok {
		return errors.New("key must be string")
	}

	dataType, exists := ctx.Storage.Type(key)
	if !exists {
		ctx.Result = "none"
	} else {
		ctx.Result = string(dataType)
	}

	return nil
}

// ExpireCommand EXPIRE 命令
type ExpireCommand struct {
	BaseCommand
}

func NewExpireCommand() *ExpireCommand {
	return &ExpireCommand{BaseCommand: BaseCommand{name: "EXPIRE"}}
}

func (c *ExpireCommand) Execute(ctx *interfaces.Context) error {
	if len(ctx.Args) < 2 {
		return errors.New("EXPIRE requires at least 2 arguments")
	}

	key, ok := ctx.Args[0].(string)
	if !ok {
		return errors.New("key must be string")
	}

	var ttl time.Duration
	switch t := ctx.Args[1].(type) {
	case time.Duration:
		ttl = t
	case int:
		ttl = time.Duration(t) * time.Second
	case int64:
		ttl = time.Duration(t) * time.Second
	default:
		return errors.New("ttl must be duration or seconds")
	}

	success := ctx.Storage.Expire(key, ttl)
	ctx.Result = success
	return nil
}

// TTLCommand TTL 命令
type TTLCommand struct {
	BaseCommand
}

func NewTTLCommand() *TTLCommand {
	return &TTLCommand{BaseCommand: BaseCommand{name: "TTL"}}
}

func (c *TTLCommand) Execute(ctx *interfaces.Context) error {
	if len(ctx.Args) < 1 {
		return errors.New("TTL requires at least 1 argument")
	}

	key, ok := ctx.Args[0].(string)
	if !ok {
		return errors.New("key must be string")
	}

	ttl, exists := ctx.Storage.TTL(key)
	if !exists {
		ctx.Result = -2 // key不存在
	} else if ttl < 0 {
		ctx.Result = -1 // 永不过期
	} else {
		ctx.Result = int(ttl.Seconds())
	}

	return nil
}

// StatsCommand STATS 命令
type StatsCommand struct {
	BaseCommand
}

func NewStatsCommand() *StatsCommand {
	return &StatsCommand{BaseCommand: BaseCommand{name: "STATS"}}
}

func (c *StatsCommand) Execute(ctx *interfaces.Context) error {
	ctx.Result = ctx.Storage.Stats()
	return nil
}
