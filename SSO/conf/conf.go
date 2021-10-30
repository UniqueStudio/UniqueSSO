package conf

import (
	"net/url"

	"github.com/UniqueStudio/UniqueSSO/pb/sso"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
	"github.com/xylonx/zapx"
	"go.uber.org/zap"
)

type Conf struct {
	Application  ApplicationConf  `mapstructure:"application"`
	Database     DatabaseConf     `mapstructure:"database"`
	Redis        RedisConf        `mapstructure:"redis"`
	Sms          []SMSOptions     `mapstructure:"sms"`
	OpenPlatform OpenPlatformConf `mapstructure:"open_platform"`
	APM          APMConf          `mapstructure:"apm"`
	Lark         LarkConf         `mapstructure:"lark"`
}
type ApplicationConf struct {
	Host             string `mapstructure:"host"`
	Port             string `mapstructure:"port"`
	RPCPort          string `mapstructure:"rpc_port"`
	Name             string `mapstructure:"name"`
	Mode             string `mapstructure:"mode"`
	ReadTimeout      int    `mapstructure:"read_timeout"`
	WriteTimeout     int    `mapstructure:"write_timeout"`
	SessionSecret    string `mapstructure:"session_secret"`
	SessionDomain    string `mapstructure:"session_domain"`
	LoginRedirectURI string `mapstructure:"login_redirect_uri"`
}

type DatabaseConf struct {
	PostgresDSN string `mapstructure:"postgres_dsn"`
}

type RedisConf struct {
	Addr     string `mapstructure:"addr"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type SMSOptions struct {
	Name       string `mapstructure:"name" validator:"oneof='verificationCode'"`
	TemplateId string `mapstructure:"template_id"`
	SignName   string `mapstructure:"sign_name"`
}

type OpenPlatformConf struct {
	GrpcAddr string `mapstructure:"grpc_addr"`
}

type APMConf struct {
	ReporterBackground string `mapstructure:"reporter_backend"`
}

type LarkConf struct {
	AppId       string `mapstructure:"app_id"`
	AppSecret   string `mapstructure:"app_secret"`
	RedirectUri string `mapstructure:"redirect_uri"`
	GroupId     struct {
		Web     string `mapstructure:"web"`
		Lab     string `mapstructure:"lab"`
		PM      string `mapstructure:"pm"`
		Design  string `mapstructure:"design"`
		Android string `mapstructure:"android"`
		IOS     string `mapstructure:"ios"`
		Game    string `mapstructure:"game"`
		AI      string `mapstructure:"ai"`
	} `mapstructure:"group_id"`
	LarkUserType struct {
		Intern    int32 `mapstructure:"intern"`
		Regular   int32 `mapstructure:"regular"`
		Graduated int32 `mapstructure:"graduated"`
	} `mapstructure:"lark_user_type"`
	LarkEnumTypeMap map[int32]sso.Role `mapstructure:"-"`
}

var (
	SSOConf = &Conf{}
)

func InitConf(confFilepath string) error {
	viper.SetConfigFile(confFilepath)
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}
	err = viper.Unmarshal(SSOConf)
	if err != nil {
		return err
	}

	validate := validator.New()
	if err := validate.Struct(SSOConf); err != nil {
		return err
	}

	if SSOConf.Application.Mode == "debug" {
		zapx.Info("run mode", zap.String("mode", SSOConf.Application.Mode))
	}

	SSOConf.Lark.RedirectUri = url.PathEscape(SSOConf.Lark.RedirectUri)

	SSOConf.Lark.LarkEnumTypeMap = make(map[int32]sso.Role)
	SSOConf.Lark.LarkEnumTypeMap[SSOConf.Lark.LarkUserType.Intern] = sso.Role_INTERN
	SSOConf.Lark.LarkEnumTypeMap[SSOConf.Lark.LarkUserType.Regular] = sso.Role_REGULAR
	SSOConf.Lark.LarkEnumTypeMap[SSOConf.Lark.LarkUserType.Graduated] = sso.Role_GRADUATED

	return nil
}
