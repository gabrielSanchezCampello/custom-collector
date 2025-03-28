# custom-collector with execreceiver
This receiver is created to allow the periodic execution of commands from the opentelemetry collector.
The result of the command to be executed must be a numerical value, which will be used as the value for the metric defined.

In this receiver it must be defined the following options for each command:
- command: the command to be executed.
- metric: in this section all the metric information to be associated to the command result will be defined.
  - metric_name: the name of the metric.
  - static_attributes: list of the attributes to be assigned to the metric-command. 

For example:

<code>
execreceiver:
    queries:
      - command:  "echo $((RANDOM % 100))"
        metric:
          metric_name: echo_ft.prueba
          static_attributes:
            region: "es"

</code>
