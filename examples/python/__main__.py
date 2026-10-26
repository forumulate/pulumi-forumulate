import pulumi
import pulumi_forumulate as forumulate

# FORUMULATE_TOKEN and FORUMULATE_ORG_ID come from the environment (or
# `pulumi config set --secret forumulate:token ...`).
config = pulumi.Config()
space_id = config.require("statusSpaceId")

groups = forumulate.StatusGroups(
    "page",
    space_id=space_id,
    groups=[
        {"name": "Core", "description": "What every customer touches."},
        {"name": "Integrations", "layout": "grid-2", "hide_uptime": True},
    ],
)

api = forumulate.StatusComponent(
    "api",
    space_id=space_id,
    name="API",
    url="https://api.example.com/health",
    group="Core",
    regions=["eu-west", "us-east"],
    confirm_regions=2,
    bar_days=90,
    interval_seconds=60,
    opts=pulumi.ResourceOptions(depends_on=[groups]),
)

backup = forumulate.StatusComponent(
    "nightly-backup",
    space_id=space_id,
    name="Nightly backup",
    kind="heartbeat",
    interval_seconds=90000,
    display="status",
)

forumulate.StatusMetric(
    "error-rate",
    space_id=space_id,
    key="error-rate",
    name="Error rate",
    unit="%",
    decimals=2,
    threshold=1,
)

pulumi.export("api_component", api.id)
pulumi.export("backup_ping_url", backup.ping_url)
