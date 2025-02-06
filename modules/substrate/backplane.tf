#
# In All Accounts
#

# Create a backplane event bus in all accounts.
resource "aws_cloudwatch_event_bus" "backplane" {
    name = "${local.resource_name}-backplane"
}

# Create a log group for the backplane event bus.
resource "aws_cloudwatch_log_group" "backplane" {
    name = "/aws/events/${aws_cloudwatch_event_bus.backplane.name}"
    retention_in_days = 1
}

# Create a logging rule for the backplane event bus.
resource "aws_cloudwatch_event_rule" "logging" {
    name = "debug-logging-route"
    event_bus_name = aws_cloudwatch_event_bus.backplane.name
    event_pattern = jsonencode({
        source = [{
            prefix = ""
        }]
    })
    state = "DISABLED"
}

# Create a logging target for the backplane event bus.
resource "aws_cloudwatch_event_target" "logging" {
    rule = aws_cloudwatch_event_rule.logging.name
    event_bus_name = aws_cloudwatch_event_bus.backplane.name
    target_id = "cloudwatch-logs"
    arn = aws_cloudwatch_log_group.backplane.arn
}

#
# In Spoke Accounts
#

# Create hub account policies for each spoke account.
resource "aws_cloudwatch_event_bus_policy" "allow_hub" {
    count = local.is_spoke ? 1 : 0
    event_bus_name = aws_cloudwatch_event_bus.backplane.name
    policy = jsonencode({
        Version = "2012-10-17"
        Statement = [
            {
                Sid = "AllowPutEventsFromHub"
                Effect = "Allow"
                Principal = {
                    AWS = "arn:aws:iam::${local.hub_account_id}:root"
                }
                Action = "events:PutEvents"
                Resource = aws_cloudwatch_event_bus.backplane.arn
            }
        ]
    })
}

#
# In Hub Account
#

# Create route rules to each spoke account.
resource "aws_cloudwatch_event_rule" "route" {
    for_each = local.is_hub ? local.spoke_map : {}
    name = "${each.value}-route"
    event_bus_name = aws_cloudwatch_event_bus.backplane.name
    event_pattern = jsonencode({
        detail = {
            header = {
                source = [local.hub_account_id, local.hub_account_name]
                destination = [each.key, each.value]
            }
        }
    })
}

# Create targets for each spoke account.
resource "aws_cloudwatch_event_target" "destination" {
    for_each = local.is_hub ? local.spoke_map : {}
    rule = aws_cloudwatch_event_rule.route[each.key].name
    event_bus_name = aws_cloudwatch_event_bus.backplane.name
    target_id = each.value
    arn       = "arn:aws:events:${data.aws_region.current.name}:${each.key}:event-bus/${local.resource_name}-backplane"
    role_arn  = aws_iam_role.router[0].arn
}

# Create a router role in the hub account.
resource "aws_iam_role" "router" {
    count = local.is_hub ? 1 : 0
    name  = "${local.resource_name}-router"
    assume_role_policy = jsonencode({
        Version = "2012-10-17"
        Statement = [
            {
                Sid = "AllowAssumeRole"
                Action = "sts:AssumeRole"
                Effect = "Allow"
                Principal = {
                    Service = "events.amazonaws.com"
                }
            }
        ]
    })
}

# Create a router policy in the hub account.
resource "aws_iam_role_policy" "router" {
    count = local.is_hub ? 1 : 0
    name  = "${local.resource_name}-router"
    role  = aws_iam_role.router[0].id
    policy = jsonencode({
        Version = "2012-10-17"
        Statement = [
            {
                Sid = "AllowPutEventsToSpokes"
                Effect = "Allow"
                Action = "events:PutEvents"
                Resource = [
                    for spoke in var.spoke_accounts:
                        "arn:aws:events:${data.aws_region.current.name}:${spoke.account_id}:event-bus/${local.resource_name}-backplane"
                ]
            }
        ]
    })
}

