import json
import logging

from opentelemetry.trace.propagation import get_current_span


def lambda_handler(event, context):
    logger = logging.getLogger()

    # Removes presumably the only handler (LambdaLoggerHandler)
    # that was set by AWS Lambda Python runtime
    all_handlers = logger.handlers
    logger.removeHandler(all_handlers[0])

    # Get current span from Otel SDK to use in formatting
    current_span = get_current_span()
    ctx = current_span.get_span_context()

    # Set up new log record factory to accept RequestId, span
    # context for formatting
    old_factory = logging.getLogRecordFactory()
    def record_factory(*args, **kwargs):
        record = old_factory(*args, **kwargs)
        record.aws_request_id = str(context.aws_request_id)
        record.otelSpanID = format(ctx.span_id, "016x")
        record.otelTraceID = format(ctx.trace_id, "032x")
        record.otelTraceSampled = ctx.trace_flags.sampled
        return record
    
    logging.setLogRecordFactory(record_factory)

    # Set up new handler with format that has RequestId, span context
    test_formatter = logging.Formatter(
        "%(asctime)s %(levelname)s [%(name)s] [%(filename)s:%(lineno)d] " \
        "[RequestId %(aws_request_id)s] " \
        "[trace_id=%(otelTraceID)s span_id=%(otelSpanID)s " \
        "trace_flags=%(otelTraceSampled)02d resource.service.name=%(otelServiceName)s] - %(message)s"
    )
    test_handler = logging.StreamHandler()
    test_handler.setFormatter(test_formatter)
    logger.addHandler(test_handler)

    # Test logs and response
    logger.warning("Here is a logger warning message")
    logger.warning("Here is context.aws_request_id: %s", context.aws_request_id)

    return {
        'statusCode': 200,
        'body': json.dumps('Hello from Lambda!')
    }
