package caas

import (
	"github.com/hashicorp/terraform-plugin-framework/function"

	"fractal.cloud/terraform-provider-fc/internal/provider/functions/spec"
)

func NewBigdataCaasDataProcessingJobFunction() function.Function {
	return spec.New(spec.Spec{
		Name:          "bigdata_caas_data_processing_job",
		ComponentType: "BigData.CaaS.DataProcessingJob",
		Summary:       "Creates a containerized Spark Data Processing Job blueprint component",
		Description: "Builds a BigData.CaaS.DataProcessingJob component (a Spark application) on a " +
			"BigData.CaaS.DistributedDataProcessing platform. job_type is required; main_class is needed for Java and Scala jobs. " +
			"Link the job to a Datalake with settings = { purpose = \"raw\" | \"curated\" | \"checkpoint\", path = ... } " +
			"and to a messaging Entity with settings = { access = \"publish\" | \"subscribe\" | \"publish-subscribe\" }.",
		Attributes: []spec.Attribute{
			{Name: "job_type", Key: "type", Kind: spec.String, Required: true, OneOf: []string{"Java", "Python", "Scala"}},
			{Name: "mode", Key: "mode", Kind: spec.String, OneOf: []string{"cluster", "client"}},
			{Name: "main_class", Key: "mainClass", Kind: spec.String},
			{Name: "main_application_file", Key: "mainApplicationFile", Kind: spec.String},
			{Name: "arguments", Key: "arguments", Kind: spec.StringList},
			{Name: "spark_conf", Key: "sparkConf", Kind: spec.StringMap},
			{Name: "schedule", Key: "schedule", Kind: spec.String},
			{Name: "restart_policy", Key: "restartPolicy", Kind: spec.String, OneOf: []string{"Never", "Always", "OnFailure"}},
			{Name: "max_retries", Key: "maxRetries", Kind: spec.Int64},
			{Name: "concurrency_policy", Key: "concurrencyPolicy", Kind: spec.String, OneOf: []string{"Allow", "Forbid", "Replace"}},
		},
		Dependencies: []spec.Dependency{{Name: "platform", ComponentType: "BigData.CaaS.DistributedDataProcessing"}},
		Links:        true,
	})
}
