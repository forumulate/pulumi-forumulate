import * as pulumi from "@pulumi/pulumi";
import * as forumulate from "@forumulate/pulumi";

// FORUMULATE_TOKEN and FORUMULATE_ORG_ID come from the environment (or
// `pulumi config set --secret forumulate:token ...`).
const config = new pulumi.Config();
const spaceId = config.require("statusSpaceId");

const groups = new forumulate.StatusGroups("page", {
    spaceId,
    groups: [
        { name: "Core", description: "What every customer touches." },
        { name: "Integrations", layout: "grid-2", hideUptime: true },
    ],
});

const api = new forumulate.StatusComponent("api", {
    spaceId,
    name: "API",
    url: "https://api.example.com/health",
    group: "Core",
    regions: ["eu-west", "us-east"],
    confirmRegions: 2,
    barDays: 90,
    intervalSeconds: 60,
}, { dependsOn: [groups] });

const backup = new forumulate.StatusComponent("nightly-backup", {
    spaceId,
    name: "Nightly backup",
    kind: "heartbeat",
    intervalSeconds: 90000,
    display: "status",
});

new forumulate.StatusMetric("error-rate", {
    spaceId,
    key: "error-rate",
    name: "Error rate",
    unit: "%",
    decimals: 2,
    threshold: 1,
});

export const apiComponent = api.id;
export const backupPingUrl = backup.pingUrl;
export const page = forumulate.getStatusPageOutput({ spaceId });
