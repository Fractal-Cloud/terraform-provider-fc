package paas

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"fractal.cloud/terraform-provider-fc/internal/provider/components"
)

var _ function.Function = &BigdataPaasDataProcessingJobFunction{}

type BigdataPaasDataProcessingJobFunction struct{}

func NewBigdataPaasDataProcessingJobFunction() function.Function {
	return &BigdataPaasDataProcessingJobFunction{}
}

func (f *BigdataPaasDataProcessingJobFunction) Metadata(_ context.Context, _ function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "bigdata_paas_data_processing_job"
}

func (f *BigdataPaasDataProcessingJobFunction) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary: "Creates a BigData PaaS Data Processing Job blueprint component",
		Description: "Builds a BigData PaaS Data Processing Job component with the correct type for use in a fractal's components list. " +
			"If platform is provided, it is automatically added as a dependency. " +
			"Link the job to a Datalake with settings = { purpose = \"raw\" | \"curated\" | \"checkpoint\", path = ... } " +
			"and to a messaging Entity with settings = { access = \"publish\" | \"subscribe\" | \"publish-subscribe\" }.",
		Parameters: []function.Parameter{
			function.ObjectParameter{
				Name:        "config",
				Description: "Data Processing Job configuration",
				CustomType: components.NewConfigObjectType(map[string]attr.Type{
					"id":               types.StringType,
					"display_name":     types.StringType,
					"description":      types.StringType,
					"platform":         components.ComponentObjectType,
					"job_name":         types.StringType,
					"task_type":        types.StringType,
					"notebook_path":    types.StringType,
					"python_file":      types.StringType,
					"main_class_name":  types.StringType,
					"jar_uri":          types.StringType,
					"cron_schedule":    types.StringType,
					"max_retries":      types.Int64Type,
					"existing_cluster": types.BoolType,
					"parameters":       types.ListType{ElemType: types.StringType},
					"links":            components.LinksAttrType,
					"extra_parameters": components.ParametersAttrType,
				}, "id"),
			},
		},
		Return: components.ComponentReturn(),
	}
}

type bigdataPaasDataProcessingJobConfig struct {
	Id              types.String  `tfsdk:"id"`
	DisplayName     types.String  `tfsdk:"display_name"`
	Description     types.String  `tfsdk:"description"`
	Platform        types.Object  `tfsdk:"platform"`
	JobName         types.String  `tfsdk:"job_name"`
	TaskType        types.String  `tfsdk:"task_type"`
	NotebookPath    types.String  `tfsdk:"notebook_path"`
	PythonFile      types.String  `tfsdk:"python_file"`
	MainClassName   types.String  `tfsdk:"main_class_name"`
	JarUri          types.String  `tfsdk:"jar_uri"`
	CronSchedule    types.String  `tfsdk:"cron_schedule"`
	MaxRetries      types.Int64   `tfsdk:"max_retries"`
	ExistingCluster types.Bool    `tfsdk:"existing_cluster"`
	Parameters      types.List    `tfsdk:"parameters"`
	Links           types.Dynamic `tfsdk:"links"`
	ExtraParameters types.Map     `tfsdk:"extra_parameters"`
}

func (f *BigdataPaasDataProcessingJobFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var config bigdataPaasDataProcessingJobConfig
	resp.Error = function.ConcatFuncErrors(resp.Error, req.Arguments.Get(ctx, &config))
	if resp.Error != nil {
		return
	}

	params := map[string]string{}

	if !config.JobName.IsNull() && !config.JobName.IsUnknown() {
		params["jobName"] = config.JobName.ValueString()
	}
	if !config.TaskType.IsNull() && !config.TaskType.IsUnknown() {
		params["taskType"] = config.TaskType.ValueString()
	}
	if !config.NotebookPath.IsNull() && !config.NotebookPath.IsUnknown() {
		params["notebookPath"] = config.NotebookPath.ValueString()
	}
	if !config.PythonFile.IsNull() && !config.PythonFile.IsUnknown() {
		params["pythonFile"] = config.PythonFile.ValueString()
	}
	if !config.MainClassName.IsNull() && !config.MainClassName.IsUnknown() {
		params["mainClassName"] = config.MainClassName.ValueString()
	}
	if !config.JarUri.IsNull() && !config.JarUri.IsUnknown() {
		params["jarUri"] = config.JarUri.ValueString()
	}
	if !config.CronSchedule.IsNull() && !config.CronSchedule.IsUnknown() {
		params["cronSchedule"] = config.CronSchedule.ValueString()
	}
	if !config.MaxRetries.IsNull() && !config.MaxRetries.IsUnknown() {
		params["maxRetries"] = fmt.Sprintf("%d", config.MaxRetries.ValueInt64())
	}
	if !config.ExistingCluster.IsNull() && !config.ExistingCluster.IsUnknown() {
		params["existingCluster"] = fmt.Sprintf("%t", config.ExistingCluster.ValueBool())
	}

	if !config.Parameters.IsNull() && !config.Parameters.IsUnknown() {
		paramList := make([]string, 0, len(config.Parameters.Elements()))
		for _, v := range config.Parameters.Elements() {
			paramList = append(paramList, v.(types.String).ValueString())
		}
		paramJSON, err := json.Marshal(paramList)
		if err != nil {
			resp.Error = function.ConcatFuncErrors(resp.Error, function.NewFuncError("failed to marshal parameters: "+err.Error()))
			return
		}
		params["parameters"] = string(paramJSON)
	}

	var deps []string
	platformId, funcErr := components.ExtractDependency(config.Platform, "BigData.PaaS.DistributedDataProcessing")
	if funcErr != nil {
		resp.Error = function.ConcatFuncErrors(resp.Error, funcErr)
		return
	}
	if platformId != "" {
		deps = append(deps, platformId)
	}

	parameters, funcErr := components.WithExtraParameters(params, config.ExtraParameters)
	if funcErr != nil {
		resp.Error = function.ConcatFuncErrors(resp.Error, funcErr)
		return
	}

	links, funcErr := components.LinksFromDynamic(config.Links)
	if funcErr != nil {
		resp.Error = function.ConcatFuncErrors(resp.Error, funcErr)
		return
	}

	result, funcErr := components.BuildComponent(
		config.Id.ValueString(),
		"BigData.PaaS.DataProcessingJob",
		components.OptionalString(config.DisplayName),
		components.OptionalString(config.Description),
		types.StringNull(),
		parameters,
		deps,
		links,
	)
	resp.Error = function.ConcatFuncErrors(resp.Error, funcErr)
	if resp.Error != nil {
		return
	}

	resp.Error = function.ConcatFuncErrors(resp.Error, resp.Result.Set(ctx, result))
}
