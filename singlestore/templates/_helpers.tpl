{{/* Define common labels */}}
{{- define "singlestore-helm.labels" -}}
app.kubernetes.io/name: {{ include "singlestore-helm.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end -}}

{{- define "singlestore-helm.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "singlestore-helm.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- printf "%s-%s" .Chart.Name .Values.clusterName | trunc 63 | trimSuffix "-" -}}
{{- end -}}
{{- end -}}

{{- define "singlestore-helm.renderNodeSpec" -}}
count: {{ .count }}
{{- if .cores }}
cores: {{ .cores }}
{{- end }}
{{- if .coresLimit }}
coresLimit: {{ .coresLimit }}
{{- end }}
{{- if .memoryMB }}
memoryMB: {{ .memoryMB }}
{{- end }}
{{- if .memoryLimitMB }}
memoryLimitMB: {{ .memoryLimitMB }}
{{- end }}
{{- if .globalVariables }}
globalVariables: {{- toYaml .globalVariables | nindent 2 }}
{{- end }}
storageGB: {{ .storageGB }}
storageClass: {{ required "valid storage class is required for every specified role" .storageClass }}
{{- if .objectMetaOverrides }}
objectMetaOverrides:
{{- if and .objectMetaOverrides.labels }}
  labels: {{- toYaml .objectMetaOverrides.labels | nindent 8 }}
{{- end }}
{{- if and .objectMetaOverrides.annotations }}
  annotations: {{- toYaml .objectMetaOverrides.annotations | nindent 8 }}
{{- end }}
{{- end }}
{{- end -}}

{{- define "singlestore-helm.renderSchedulingDetailsForRole" -}}
{{- if .nodeSelector }}
nodeSelector: {{- toYaml .nodeSelector | nindent 2 }}
{{- end }}
{{- if .tolerations }}
tolerations: {{- toYaml .tolerations | nindent 2 }}
{{- end }}
{{- end -}}

{{- define "singlestore-helm.exporterHost" -}}
{{- if .Values.monitoringJob.exporterHost }}
{{- .Values.monitoringJob.exporterHost }}
{{- else }}
{{- printf "node-%s-master-0.svc-%s.%s.svc.cluster.local" .Values.clusterName .Values.clusterName .Release.Namespace }}
{{- end }}
{{- end -}}