package common

import (
	"strings"
)

const (
	// NoneProfile Profile
	NoneProfile Profile = "none"
	// DevelopPorfile Profile, 开发环境
	DevelopPorfile Profile = "develop"
	// TestProfile  Profile, 测试环境
	TestProfile Profile = "test"
	// ProductionProfile Profile, 线上环境
	ProductionProfile Profile = "production"
)

func (profile Profile) String() string {
	switch profile {
	case DevelopPorfile:
		return "develop"
	case TestProfile:
		return "test"
	case ProductionProfile:
		return "production"
	default:
		return "none"
	}
}

// ParseProfile 字符串转换成Profile类型
func ParseProfile(profileName string) Profile {
	switch strings.ToLower(profileName) {
	case "develop":
		return DevelopPorfile
	case "test":
		return TestProfile
	case "production":
		return ProductionProfile
	}
	return NoneProfile
}
