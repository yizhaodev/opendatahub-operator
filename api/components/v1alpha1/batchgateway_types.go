package v1alpha1

import (
	operatorv1 "github.com/openshift/api/operator/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/opendatahub-io/opendatahub-operator/v2/api/common"
)

const (
	BatchGatewayComponentName = "batchgateway"
	BatchGatewayInstanceName  = "default-" + BatchGatewayComponentName
	BatchGatewayKind          = "LLMBatchGateway"
)

var _ common.PlatformObject = (*LLMBatchGateway)(nil)

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster
// +kubebuilder:validation:XValidation:rule="self.metadata.name == 'default-batchgateway'",message="LLMBatchGateway name must be default-batchgateway"
// +kubebuilder:printcolumn:name="Ready",type=string,JSONPath=`.status.conditions[?(@.type=="Ready")].status`,description="Ready"
// +kubebuilder:printcolumn:name="Reason",type=string,JSONPath=`.status.conditions[?(@.type=="Ready")].reason`,description="Reason"

type LLMBatchGateway struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   LLMBatchGatewaySpec   `json:"spec,omitempty"`
	Status LLMBatchGatewayStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true
type LLMBatchGatewayList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []LLMBatchGateway `json:"items"`
}

type LLMBatchGatewaySpec struct {
	// +kubebuilder:validation:Required
	SecretRef BatchGatewaySecretReference `json:"secretRef"`

	// +kubebuilder:validation:Enum=redis;postgresql;valkey
	// +kubebuilder:default=postgresql
	DBBackend string `json:"dbBackend,omitempty"`

	FileStorage *BatchGatewayFileStorageSpec `json:"fileStorage,omitempty"`

	// +kubebuilder:validation:Required
	APIServer BatchGatewayAPIServerSpec `json:"apiServer"`

	// +kubebuilder:validation:Required
	Processor BatchGatewayProcessorSpec `json:"processor"`

	GC BatchGatewayGCSpec `json:"gc"`

	Monitoring     *BatchGatewayMonitoringSpec     `json:"monitoring,omitempty"`
	Grafana        *BatchGatewayGrafanaSpec        `json:"grafana,omitempty"`
	TLS            *BatchGatewayTLSSpec            `json:"tls,omitempty"`
	HTTPRoute      *BatchGatewayHTTPRouteSpec      `json:"httpRoute,omitempty"`
	OTEL           *BatchGatewayOTELSpec           `json:"otel,omitempty"`
	PrometheusRule *BatchGatewayPrometheusRuleSpec `json:"prometheusRule,omitempty"`
}

type BatchGatewaySecretReference struct {
	// +kubebuilder:validation:Required
	Name string `json:"name"`
}

// --- File Storage ---

type BatchGatewayFileStorageSpec struct {
	// +kubebuilder:validation:Enum=fs;s3
	// +kubebuilder:default=s3
	Type string `json:"type,omitempty"`

	S3    *BatchGatewayS3StorageSpec `json:"s3,omitempty"`
	FS    *BatchGatewayFSStorageSpec `json:"fs,omitempty"`
	Retry *BatchGatewayFileRetrySpec `json:"retry,omitempty"`
}

type BatchGatewayS3StorageSpec struct {
	Region           string `json:"region,omitempty"`
	Endpoint         string `json:"endpoint,omitempty"`
	AccessKeyID      string `json:"accessKeyId,omitempty"`
	Prefix           string `json:"prefix,omitempty"`
	UsePathStyle     bool   `json:"usePathStyle,omitempty"`
	AutoCreateBucket bool   `json:"autoCreateBucket,omitempty"`
}

type BatchGatewayFSStorageSpec struct {
	BasePath string `json:"basePath,omitempty"`
	PVCName  string `json:"pvcName,omitempty"`
}

type BatchGatewayFileRetrySpec struct {
	// +kubebuilder:default=3
	MaxRetries int32 `json:"maxRetries,omitempty"`

	// +kubebuilder:default="1s"
	InitialBackoff string `json:"initialBackoff,omitempty"`

	// +kubebuilder:default="10s"
	MaxBackoff string `json:"maxBackoff,omitempty"`
}

// --- API Server ---

type BatchGatewayAPIServerSpec struct {
	// +kubebuilder:default=1
	// +kubebuilder:validation:Minimum=0
	Replicas *int32 `json:"replicas,omitempty"`

	// +kubebuilder:validation:Required
	Image string `json:"image"`

	Resources *corev1.ResourceRequirements `json:"resources,omitempty"`

	Config *BatchGatewayAPIServerConfigSpec `json:"config,omitempty"`
}

type BatchGatewayAPIServerConfigSpec struct {
	Port                string `json:"port,omitempty"`
	ObservabilityPort   string `json:"observabilityPort,omitempty"`
	ReadTimeoutSeconds  int32  `json:"readTimeoutSeconds,omitempty"`
	WriteTimeoutSeconds int32  `json:"writeTimeoutSeconds,omitempty"`
	IdleTimeoutSeconds  int32  `json:"idleTimeoutSeconds,omitempty"`

	BatchAPI *BatchGatewayBatchAPIConfig `json:"batchAPI,omitempty"`
	FileAPI  *BatchGatewayFileAPIConfig  `json:"fileAPI,omitempty"`

	EnablePprof bool                       `json:"enablePprof,omitempty"`
	Logging     *BatchGatewayLoggingConfig `json:"logging,omitempty"`
}

type BatchGatewayBatchAPIConfig struct {
	EventTTLSeconds    int32    `json:"eventTTLSeconds,omitempty"`
	PassThroughHeaders []string `json:"passThroughHeaders,omitempty"`
}

type BatchGatewayFileAPIConfig struct {
	DefaultExpirationSeconds int64 `json:"defaultExpirationSeconds,omitempty"`
	MaxSizeBytes             int64 `json:"maxSizeBytes,omitempty"`
	MaxLineCount             int64 `json:"maxLineCount,omitempty"`
}

type BatchGatewayLoggingConfig struct {
	Verbosity int32 `json:"verbosity,omitempty"`
}

// --- Processor ---

type BatchGatewayProcessorSpec struct {
	// +kubebuilder:default=1
	// +kubebuilder:validation:Minimum=0
	Replicas *int32 `json:"replicas,omitempty"`

	// +kubebuilder:validation:Required
	Image string `json:"image"`

	Resources *corev1.ResourceRequirements `json:"resources,omitempty"`

	GlobalInferenceGateway *BatchGatewayInferenceGatewaySpec           `json:"globalInferenceGateway,omitempty"`
	ModelGateways          map[string]BatchGatewayInferenceGatewaySpec `json:"modelGateways,omitempty"`

	Config *BatchGatewayProcessorConfigSpec `json:"config,omitempty"`
}

type BatchGatewayInferenceGatewaySpec struct {
	// +kubebuilder:validation:Required
	URL string `json:"url"`

	RequestTimeout string `json:"requestTimeout,omitempty"`
	MaxRetries     *int32 `json:"maxRetries,omitempty"`
	InitialBackoff string `json:"initialBackoff,omitempty"`
	MaxBackoff     string `json:"maxBackoff,omitempty"`

	TLSInsecureSkipVerify bool   `json:"tlsInsecureSkipVerify,omitempty"`
	TLSCACertFile         string `json:"tlsCaCertFile,omitempty"`
	TLSClientCertFile     string `json:"tlsClientCertFile,omitempty"`
	TLSClientKeyFile      string `json:"tlsClientKeyFile,omitempty"`
}

type BatchGatewayProcessorConfigSpec struct {
	NumWorkers             int32  `json:"numWorkers,omitempty"`
	GlobalConcurrency      int32  `json:"globalConcurrency,omitempty"`
	PerModelMaxConcurrency int32  `json:"perModelMaxConcurrency,omitempty"`
	RecoveryMaxConcurrency int32  `json:"recoveryMaxConcurrency,omitempty"`
	InferenceObjective     string `json:"inferenceObjective,omitempty"`

	DefaultOutputExpirationSeconds int64 `json:"defaultOutputExpirationSeconds,omitempty"`
	ProgressTTLSeconds             int64 `json:"progressTTLSeconds,omitempty"`

	EnablePprof bool                       `json:"enablePprof,omitempty"`
	Logging     *BatchGatewayLoggingConfig `json:"logging,omitempty"`
}

// --- GC ---

type BatchGatewayGCSpec struct {
	// +kubebuilder:validation:Required
	Image string `json:"image"`

	// +kubebuilder:default="30m"
	Interval string `json:"interval,omitempty"`

	Config *BatchGatewayGCConfigSpec `json:"config,omitempty"`
}

type BatchGatewayGCConfigSpec struct {
	DryRun         bool  `json:"dryRun,omitempty"`
	MaxConcurrency int32 `json:"maxConcurrency,omitempty"`

	Logging *BatchGatewayLoggingConfig `json:"logging,omitempty"`
}

// --- Observability ---

type BatchGatewayMonitoringSpec struct {
	Enabled bool `json:"enabled,omitempty"`
}

type BatchGatewayGrafanaSpec struct {
	Enabled bool `json:"enabled,omitempty"`
}

type BatchGatewayPrometheusRuleSpec struct {
	Enabled bool              `json:"enabled,omitempty"`
	Labels  map[string]string `json:"labels,omitempty"`
}

type BatchGatewayOTELSpec struct {
	Endpoint          string `json:"endpoint,omitempty"`
	Insecure          bool   `json:"insecure,omitempty"`
	Sampler           string `json:"sampler,omitempty"`
	SamplerArg        string `json:"samplerArg,omitempty"`
	RedisTracing      bool   `json:"redisTracing,omitempty"`
	PostgresqlTracing bool   `json:"postgresqlTracing,omitempty"`
}

// --- TLS ---

type BatchGatewayTLSSpec struct {
	Enabled     bool                         `json:"enabled,omitempty"`
	SecretName  string                       `json:"secretName,omitempty"`
	CertManager *BatchGatewayCertManagerSpec `json:"certManager,omitempty"`
}

type BatchGatewayCertManagerSpec struct {
	IssuerName string   `json:"issuerName,omitempty"`
	IssuerKind string   `json:"issuerKind,omitempty"`
	DNSNames   []string `json:"dnsNames,omitempty"`
}

// --- HTTPRoute ---

type BatchGatewayHTTPRouteSpec struct {
	Enabled     bool                          `json:"enabled,omitempty"`
	Annotations map[string]string             `json:"annotations,omitempty"`
	ParentRefs  []BatchGatewayParentReference `json:"parentRefs,omitempty"`
}

type BatchGatewayParentReference struct {
	Name        string `json:"name"`
	Namespace   string `json:"namespace,omitempty"`
	SectionName string `json:"sectionName,omitempty"`
}

// --- Status ---

type LLMBatchGatewayStatus struct {
	common.Status `json:",inline"`
}

// --- DSC integration ---

// DSCBatchGatewaySpec enables BatchGateway integration under Kserve
type DSCBatchGatewaySpec struct {
	// +kubebuilder:validation:Enum=Managed;Removed
	// +kubebuilder:default=Removed
	ManagementState operatorv1.ManagementState `json:"managementState,omitempty"`
}

// DSCBatchGatewayStatus contains the observed state of the BatchGateway exposed in the DSC instance
type DSCBatchGatewayStatus struct {
	common.ManagementSpec `json:",inline"`
}

func init() {
	SchemeBuilder.Register(&LLMBatchGateway{}, &LLMBatchGatewayList{})
}

func (c *LLMBatchGateway) GetStatus() *common.Status {
	return &c.Status.Status
}

func (c *LLMBatchGateway) GetConditions() []common.Condition {
	return c.Status.GetConditions()
}

func (c *LLMBatchGateway) SetConditions(conditions []common.Condition) {
	c.Status.SetConditions(conditions)
}
