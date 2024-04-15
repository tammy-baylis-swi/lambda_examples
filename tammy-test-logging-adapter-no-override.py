"""
DOES NOT ADD CONTEXT TO LOGS.

Works with APM/OTel Python layer when LOG_CORRELATION not enabled
and OTEL_PYTHON_DISABLED_INSTRUMENTATIONS: logging
"""

import json
import logging
import os

from opentelemetry.trace import get_tracer_provider
from opentelemetry.trace.propagation import get_current_span


def calculate_service_name():
    """Returns service name from OTel resource else use AWS_LAMBDA_FUNCTION_NAME"""
    provider = get_tracer_provider()
    resource = getattr(provider, "resource", None)
    if resource:
        otel_service_name = resource.attributes.get("service.name", None)
        if otel_service_name and otel_service_name != "unknown_service":
            return otel_service_name
    return os.environ.get("AWS_LAMBDA_FUNCTION_NAME")
    

def lambda_handler(event, context):
    # Get current span from Otel SDK to use in formatting
    current_span = get_current_span()
    ctx = current_span.get_span_context()
    
    # Get service name: resource else AWS_LAMBDA_FUNCTION_NAME
    service_name = calculate_service_name()

    logger = logging.getLogger()
    my_log_adapter = logging.LoggerAdapter(
        logger,
        {
            "otelSpanID": format(ctx.span_id, "016x"),
            "otelTraceID": format(ctx.trace_id, "032x"),
            "otelTraceSampled": ctx.trace_flags.sampled,
            "otelServiceName": service_name,
        },
    )
    
    # Test logs and response
    my_log_adapter.warning("Here is a logger warning message")
    my_log_adapter.warning("Here is context.aws_request_id: %s", context.aws_request_id)


    return {
        'statusCode': 200,
        'body': json.dumps('Hello from Lambda!')
    }
