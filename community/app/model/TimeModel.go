package model

import "time"

var GarbageCacheTime = 5 * time.Minute
var CacheTime = 10 * time.Minute
var FreshTokenEffectiveTime = 7 * 24 * time.Hour
var AccessTokenEffectiveTime = 1 * time.Minute
var TimeOutLimit = 5 * time.Second
