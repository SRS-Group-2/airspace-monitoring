resource "google_monitoring_dashboard" "firestore_dashboard" {
  project        = var.project_id
  dashboard_json = <<EOF
{
  "displayName": "Firebase",
  "mosaicLayout": {
    "columns": 12,
    "tiles": [
      {
        "height": 4,
        "widget": {
          "title": "Document Reads [SUM]",
          "xyChart": {
            "chartOptions": {
              "mode": "COLOR"
            },
            "dataSets": [
              {
                "minAlignmentPeriod": "60s",
                "plotType": "LINE",
                "targetAxis": "Y1",
                "timeSeriesQuery": {
                  "timeSeriesFilter": {
                    "aggregation": {
                      "alignmentPeriod": "60s",
                      "perSeriesAligner": "ALIGN_SUM"
                    },
                    "filter": "metric.type=\"firestore.googleapis.com/document/read_count\" resource.type=\"firestore_instance\"",
                    "secondaryAggregation": {
                      "alignmentPeriod": "60s"
                    }
                  }
                }
              }
            ],
            "timeshiftDuration": "0s",
            "yAxis": {
              "label": "y1Axis",
              "scale": "LINEAR"
            }
          }
        },
        "width": 6
      },
      {
        "height": 4,
        "widget": {
          "title": "Document Writes [SUM]",
          "xyChart": {
            "chartOptions": {
              "mode": "COLOR"
            },
            "dataSets": [
              {
                "minAlignmentPeriod": "60s",
                "plotType": "LINE",
                "targetAxis": "Y1",
                "timeSeriesQuery": {
                  "timeSeriesFilter": {
                    "aggregation": {
                      "alignmentPeriod": "60s",
                      "perSeriesAligner": "ALIGN_SUM"
                    },
                    "filter": "metric.type=\"firestore.googleapis.com/document/write_count\" resource.type=\"firestore_instance\"",
                    "secondaryAggregation": {
                      "alignmentPeriod": "60s"
                    }
                  }
                }
              }
            ],
            "timeshiftDuration": "0s",
            "yAxis": {
              "label": "y1Axis",
              "scale": "LINEAR"
            }
          }
        },
        "width": 6,
        "xPos": 6
      },
      {
        "height": 4,
        "widget": {
          "title": "Document Deletes [SUM]",
          "xyChart": {
            "chartOptions": {
              "mode": "COLOR"
            },
            "dataSets": [
              {
                "minAlignmentPeriod": "60s",
                "plotType": "LINE",
                "targetAxis": "Y1",
                "timeSeriesQuery": {
                  "timeSeriesFilter": {
                    "aggregation": {
                      "alignmentPeriod": "60s",
                      "perSeriesAligner": "ALIGN_SUM"
                    },
                    "filter": "metric.type=\"firestore.googleapis.com/document/delete_count\" resource.type=\"firestore_instance\"",
                    "secondaryAggregation": {
                      "alignmentPeriod": "60s"
                    }
                  }
                }
              }
            ],
            "timeshiftDuration": "0s",
            "yAxis": {
              "label": "y1Axis",
              "scale": "LINEAR"
            }
          }
        },
        "width": 6,
        "yPos": 4
      }
    ]
  },
  "name": "projects/477074905358/dashboards/648af953-470d-46c2-b1ff-2d1b9765e41c"
}

EOF
}
