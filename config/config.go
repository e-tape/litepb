package config

import litepb "github.com/e-tape/litepb/proto"

type Config struct {
	SourceRelative    bool            `yaml:"source_relative"`
	MemPoolMessageAll litepb.Activity `yaml:"mem_pool_message_all"`
	MemPoolListAll    litepb.Activity `yaml:"mem_pool_list_all"`
	MemPoolMapAll     litepb.Activity `yaml:"mem_pool_map_all"`
	MemPoolOneofAll   litepb.Activity `yaml:"mem_pool_oneof_all"`
}
