package provider

import (
	"context"
	"os"

	"fractal.cloud/terraform-provider-fc/internal/client"
	api_caas "fractal.cloud/terraform-provider-fc/internal/provider/functions/api_management/caas"
	api_paas "fractal.cloud/terraform-provider-fc/internal/provider/functions/api_management/paas"
	api_saas "fractal.cloud/terraform-provider-fc/internal/provider/functions/api_management/saas"
	bd_paas "fractal.cloud/terraform-provider-fc/internal/provider/functions/bigdata/paas"
	bd_saas "fractal.cloud/terraform-provider-fc/internal/provider/functions/bigdata/saas"
	cw_caas "fractal.cloud/terraform-provider-fc/internal/provider/functions/custom_workloads/caas"
	cw_faas "fractal.cloud/terraform-provider-fc/internal/provider/functions/custom_workloads/faas"
	cw_iaas "fractal.cloud/terraform-provider-fc/internal/provider/functions/custom_workloads/iaas"
	cw_paas "fractal.cloud/terraform-provider-fc/internal/provider/functions/custom_workloads/paas"
	cw_saas "fractal.cloud/terraform-provider-fc/internal/provider/functions/custom_workloads/saas"
	msg_caas "fractal.cloud/terraform-provider-fc/internal/provider/functions/messaging/caas"
	msg_paas "fractal.cloud/terraform-provider-fc/internal/provider/functions/messaging/paas"
	msg_saas "fractal.cloud/terraform-provider-fc/internal/provider/functions/messaging/saas"
	nc_iaas "fractal.cloud/terraform-provider-fc/internal/provider/functions/network_and_compute/iaas"
	nc_paas "fractal.cloud/terraform-provider-fc/internal/provider/functions/network_and_compute/paas"
	nc_saas "fractal.cloud/terraform-provider-fc/internal/provider/functions/network_and_compute/saas"
	obs_caas "fractal.cloud/terraform-provider-fc/internal/provider/functions/observability/caas"
	obs_saas "fractal.cloud/terraform-provider-fc/internal/provider/functions/observability/saas"
	"fractal.cloud/terraform-provider-fc/internal/provider/functions/secrets"
	sec_caas "fractal.cloud/terraform-provider-fc/internal/provider/functions/security/caas"
	sec_saas "fractal.cloud/terraform-provider-fc/internal/provider/functions/security/saas"
	st_caas "fractal.cloud/terraform-provider-fc/internal/provider/functions/storage/caas"
	st_paas "fractal.cloud/terraform-provider-fc/internal/provider/functions/storage/paas"
	st_saas "fractal.cloud/terraform-provider-fc/internal/provider/functions/storage/saas"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ provider.Provider              = &fractalCloudProvider{}
	_ provider.ProviderWithFunctions = &fractalCloudProvider{}
)

// New is a helper function to simplify provider server and testing implementation.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &fractalCloudProvider{
			version: version,
		}
	}
}

// fractalCloudProvider is the provider implementation.
type fractalCloudProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// fractalCloudProviderModel maps provider schema data to a Go type.
type fractalCloudProviderModel struct {
	Host                 types.String `tfsdk:"host"`
	ServiceAccountId     types.String `tfsdk:"service_account_id"`
	ServiceAccountSecret types.String `tfsdk:"service_account_secret"`
}

// Metadata returns the provider type name.
func (p *fractalCloudProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "fc"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *fractalCloudProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				Optional: true,
			},
			"service_account_id": schema.StringAttribute{
				Optional: true,
			},
			"service_account_secret": schema.StringAttribute{
				Optional: true,
			},
		},
	}
}

func (p *fractalCloudProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	// Retrieve provider data from configuration
	var config fractalCloudProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.

	host := "https://api.fractal.cloud"
	serviceAccountId := os.Getenv("FRACTAL_CLOUD_SERVICE_ACCOUNT_ID")
	serviceAccountSecret := os.Getenv("FRACTAL_CLOUD_SERVICE_ACCOUNT_SECRET")

	if !config.Host.IsNull() {
		host = config.Host.ValueString()
	}

	if !config.ServiceAccountId.IsNull() {
		serviceAccountId = config.ServiceAccountId.ValueString()
	}

	if !config.ServiceAccountSecret.IsNull() {
		serviceAccountSecret = config.ServiceAccountSecret.ValueString()
	}

	if serviceAccountId == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("service_account_id"),
			"Missing Fractal Cloud API Service Account Id",
			"The provider cannot create the Fractal Cloud API client as there is a missing or empty value for the Fractal Cloud API Service Account Id. "+
				"Set the username value in the configuration or use the FRACTAL_CLOUD_SERVICE_ACCOUNT_ID environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if serviceAccountSecret == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("service_account_secret"),
			"Missing Fractal Cloud API Service Account Secret",
			"The provider cannot create the Fractal Cloud API client as there is a missing or empty value for the Fractal Cloud API Service Account Secret. "+
				"Set the service_account_secret value in the configuration or use the FRACTAL_CLOUD_SERVICE_ACCOUNT_SECRET environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Create a new Fractal Cloud client using the configuration values
	client := fractalCloud.NewClient(&fractalCloud.ClientLogger{
		Debug: func(s string) {
			tflog.Debug(ctx, s)
		},
		Information: func(s string) {
			tflog.Info(ctx, s)
		},
		Warning: func(s string) {
			tflog.Warn(ctx, s)
		},
		Error: func(s string) {
			tflog.Error(ctx, s)
		},
	}, &host, &serviceAccountId, &serviceAccountSecret)

	// Make the Fractal Cloud client available during DataSource and Resource
	// type Configure methods.
	resp.DataSourceData = client
	resp.ResourceData = client
}

// DataSources defines the data sources implemented in the provider.
func (p *fractalCloudProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewPersonalBoundedContextDataSource,
		NewOrganizationalBoundedContextDataSource,
		NewOrganizationDataSource,
		NewFractalDataSource,
	}
}

// Resources defines the resources implemented in the provider.
func (p *fractalCloudProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewPersonalBoundedContext,
		NewOrganizationalBoundedContext,
		NewManagementEnvironment,
		NewOperationalEnvironment,
		NewFractal,
	}
}

// Functions defines the provider functions for building blueprint components.
func (p *fractalCloudProvider) Functions(_ context.Context) []func() function.Function {
	return []func() function.Function{
		// Values
		secrets.NewSecretRefFunction,

		// NetworkAndCompute
		nc_iaas.NewVirtualNetworkFunction,
		nc_iaas.NewSubnetFunction,
		nc_iaas.NewLoadBalancerFunction,
		nc_iaas.NewSecurityGroupFunction,
		nc_iaas.NewVirtualMachineFunction,
		nc_paas.NewContainerPlatformFunction,
		nc_saas.NewUnmanagedFunction,

		// CustomWorkloads
		cw_caas.NewWorkloadFunction,
		cw_iaas.NewWorkloadFunction,
		cw_paas.NewWorkloadFunction,
		cw_faas.NewWorkloadFunction,
		cw_saas.NewUnmanagedFunction,

		// Storage
		st_paas.NewStoragePaasFilesAndBlobsFunction,
		st_paas.NewStoragePaasRelationalDbmsFunction,
		st_paas.NewStoragePaasRelationalDatabaseFunction,
		st_paas.NewStoragePaasDocumentDbmsFunction,
		st_paas.NewStoragePaasDocumentDatabaseFunction,
		st_paas.NewStoragePaasColumnOrientedDbmsFunction,
		st_paas.NewStoragePaasColumnOrientedEntityFunction,
		st_paas.NewStoragePaasKeyValueDbmsFunction,
		st_paas.NewStoragePaasKeyValueEntityFunction,
		st_paas.NewStoragePaasGraphDbmsFunction,
		st_paas.NewStoragePaasGraphDatabaseFunction,
		st_caas.NewStorageCaasSearchFunction,
		st_caas.NewStorageCaasSearchEntityFunction,
		st_saas.NewStorageSaasUnmanagedFunction,

		// Messaging
		msg_paas.NewMessagingPaasBrokerFunction,
		msg_paas.NewMessagingPaasEntityFunction,
		msg_caas.NewMessagingCaasBrokerFunction,
		msg_caas.NewMessagingCaasEntityFunction,
		msg_saas.NewMessagingSaasUnmanagedFunction,

		// BigData
		bd_paas.NewBigdataPaasDistributedDataProcessingFunction,
		bd_paas.NewBigdataPaasComputeClusterFunction,
		bd_paas.NewBigdataPaasDataProcessingJobFunction,
		bd_paas.NewBigdataPaasMlExperimentFunction,
		bd_paas.NewBigdataPaasDatalakeFunction,
		bd_saas.NewBigdataSaasUnmanagedFunction,

		// APIManagement
		api_paas.NewPaaSAPIGatewayFunction,
		api_caas.NewCaaSAPIGatewayFunction,
		api_saas.NewSaaSUnmanagedFunction,

		// Observability
		obs_caas.NewCaaSMonitoringFunction,
		obs_caas.NewCaaSTracingFunction,
		obs_caas.NewCaaSLoggingFunction,
		obs_saas.NewSaaSUnmanagedFunction,

		// Security
		sec_caas.NewCaaSServiceMeshSecurityFunction,
		sec_saas.NewSaaSUnmanagedFunction,
	}
}
