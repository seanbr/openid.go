package openid

import "github.com/garyburd/redigo/redis"
import "log"

type DiscoveredInfo interface {
	OpEndpoint() string
	OpLocalId() string
	ClaimedId() string
	// ProtocolVersion: it's always openId 2.
}

type DiscoveryCache interface {
	Put(id string, info DiscoveredInfo)
	// Return a discovered info, or nil.
	Get(id string) DiscoveredInfo
}

//type SimpleDiscoveredInfo struct {
//opEndpoint string
//opLocalId  string
//claimedId  string
//}

//func (s *SimpleDiscoveredInfo) OpEndpoint() string {
//return s.opEndpoint
//}

//func (s *SimpleDiscoveredInfo) OpLocalId() string {
//return s.opLocalId
//}

//func (s *SimpleDiscoveredInfo) ClaimedId() string {
//return s.claimedId
//}

type RedisDiscoveryCache struct {
	Pool *redis.Pool
}

type RedisDiscoveredInfo struct {
	opEndpoint string
	opLocalId  string
	claimedId  string
}

func (s *RedisDiscoveredInfo) OpEndpoint() string {
	return s.opEndpoint
}

func (s *RedisDiscoveredInfo) OpLocalId() string {
	return s.opLocalId
}

func (s *RedisDiscoveredInfo) ClaimedId() string {
	return s.claimedId
}

func (s RedisDiscoveryCache) Put(id string, info DiscoveredInfo) {
	log.Print("--------------PUT--------------")
}

func (s RedisDiscoveryCache) Get(id string) DiscoveredInfo {
	info, err := s.Pool.Get().Do("GET", "steam_discovery_cache:"+id)
	if info != nil {
		log.Print(info)
		//return info
	}
	log.Print(err)
	return nil
}

func compareDiscoveredInfo(a DiscoveredInfo, opEndpoint, opLocalId, claimedId string) bool {
	return a != nil &&
		a.OpEndpoint() == opEndpoint &&
		a.OpLocalId() == opLocalId &&
		a.ClaimedId() == claimedId
}
