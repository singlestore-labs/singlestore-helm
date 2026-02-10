package template

import "fmt"

const (
	relativeChartPath    = "../../singlestore"
	clusterTemplateName  = "cluster"
	operatorTemplateName = "operator-deployment"
	sysreqDaemonSetName  = "sysreq-daemonset"
	monitoringJobName    = "monitoring-job"
)

var disabledByDefaultTemplates = []string{
	sysreqDaemonSetName,
	monitoringJobName,
}

func getTemplatePath(templateName string) string {
	return fmt.Sprintf("templates/%s.yaml", templateName)
}
