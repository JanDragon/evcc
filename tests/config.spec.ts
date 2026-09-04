import { test, expect } from "@playwright/test";
import { start, stop, restart, baseUrl } from "./evcc";
import { expectModalHidden, expectModalVisible, openMoreMenu } from "./utils";

const CONFIG_GRID_ONLY = "config-grid-only.evcc.yaml";
const NETWORK_HOST = "somehostname.local";
const NETWORK_EXTERNAL_URL = "https://ext.evcc.example";

test.use({ baseURL: baseUrl() });

test.beforeAll(async () => {
  await start(CONFIG_GRID_ONLY);
});
test.afterAll(async () => {
  await stop();
});

test.describe("basics", async () => {
  test("navigation to config", async ({ page }) => {
    await page.goto("/");
    await openMoreMenu(page);
    await page.getByRole("link", { name: "Configuration" }).click();
    await expect(page.getByRole("heading", { name: "Configuration" })).toBeVisible();
  });
});

test.describe("general", async () => {
  test("change site title", async ({ page }) => {
    // initial value on main ui
    await page.goto("/");
    await expect(page.getByRole("heading", { name: "Hello World" })).toBeVisible();

    // change value in config
    await page.goto("/#/config");

    await expect(page.getByTestId("generalconfig-title")).toContainText("Hello World");
    await page.getByTestId("generalconfig-title").getByRole("button", { name: "edit" }).click();
    const modal = page.getByTestId("title-modal");
    await expectModalVisible(modal);
    await modal.getByLabel("Title").fill("Whoops World");

    // close modal and ignore entry on cancel
    await modal.getByRole("button", { name: "Cancel" }).click();
    await expectModalHidden(modal);
    await expect(page.getByTestId("generalconfig-title")).toContainText("Hello World");

    // change and save value
    await page.getByTestId("generalconfig-title").getByRole("button", { name: "edit" }).click();
    await expectModalVisible(modal);
    await modal.getByLabel("Title").fill("Ahoy World");
    await modal.getByRole("button", { name: "Save" }).click();
    await expectModalHidden(modal);
    await expect(page.getByTestId("generalconfig-title")).toContainText("Ahoy World");

    // check changed value on main ui
    await page.getByRole("link", { name: "Charge" }).click();
    await expect(page.getByRole("heading", { name: "Ahoy World" })).toBeVisible();
  });
  test("enable experimental", async ({ page }) => {
    await page.goto("/#/config");

    const experimentalEntry = page.getByTestId("generalconfig-experimental");

    await expect(experimentalEntry).toContainText("off");
    await experimentalEntry.getByRole("button", { name: "edit" }).click();

    const modal = page.getByTestId("experimental-modal");
    await expectModalVisible(modal);

    const experimentalInput = modal.getByLabel("Enable experimental features.");
    await expect(experimentalInput).not.toBeChecked();
    await experimentalInput.click();
    await modal.getByRole("button", { name: "Close" }).click();
    await expectModalHidden(modal);
    await expect(experimentalEntry).toContainText("on");

    await restart(CONFIG_GRID_ONLY);
    await page.reload();

    await expect(experimentalEntry).toContainText("on");
    await experimentalEntry.getByRole("button", { name: "edit" }).click();
    await expectModalVisible(modal);
    await expect(experimentalInput).toBeChecked();
  });
});

test.describe("network modal", async () => {
  test("persists host and external url across restart", async ({ page }) => {
    await page.goto("/#/config");

    const networkEntry = page.getByTestId("generalconfig-network");
    await expect(networkEntry).toBeVisible();
    await networkEntry.getByRole("button", { name: "edit" }).click();

    const modal = page.getByTestId("network-modal");
    await expectModalVisible(modal);

    const portValue = await modal.getByLabel("Port").inputValue();
    await modal.getByLabel("External URL").fill(NETWORK_EXTERNAL_URL + "/somepath");
    await expect(modal.getByRole("status")).toContainText("doesn't need a path");
    await modal.getByLabel("External URL").fill(NETWORK_EXTERNAL_URL);
    await expect(modal.getByRole("status")).not.toBeVisible();
    await modal.getByLabel("mDNS Hostname").fill(NETWORK_HOST);

    await modal.getByRole("button", { name: "Save" }).click();
    await expectModalHidden(modal);

    // values immediatelly visible
    await expect(networkEntry).toContainText(portValue);

    // check persistance
    await restart(CONFIG_GRID_ONLY);
    await page.reload();

    await expect(networkEntry).toBeVisible();
    await networkEntry.getByRole("button", { name: "edit" }).click();
    await expectModalVisible(modal);

    await expect(modal.getByLabel("mDNS Hostname")).toHaveValue(NETWORK_HOST);
    await expect(modal.getByLabel("Port")).toHaveValue(portValue);
    await expect(modal.getByLabel("External URL")).toHaveValue(NETWORK_EXTERNAL_URL);
    const internalUrl = await modal.getByLabel("Internal URL").inputValue();
    expect(internalUrl).toContain(`:${portValue}`);

    await modal.getByRole("button", { name: "Cancel" }).click();
    await expectModalHidden(modal);
  });
});

test.describe("control modal", async () => {
  test("persists operating voltage and residual power across restart", async ({ page }) => {
    await page.goto("/#/config");

    const controlEntry = page.getByTestId("generalconfig-control");
    await expect(controlEntry).toBeVisible();
    await controlEntry.getByRole("button", { name: "edit" }).click();

    const modal = page.getByTestId("control-modal");
    await expectModalVisible(modal);

    await expect(modal.getByLabel("Operating voltage")).toHaveValue("230");
    await expect(modal.getByLabel("Residual power")).toHaveValue("0");

    await modal.getByLabel("Update interval").fill("10");
    await modal.getByLabel("Operating voltage").fill("240");
    await modal.getByLabel("Residual power").fill("150");
    await modal.getByRole("button", { name: "Save" }).click();
    await expectModalHidden(modal);

    await expect(controlEntry).toContainText("10s");

    await restart(CONFIG_GRID_ONLY);
    await page.reload();

    await controlEntry.getByRole("button", { name: "edit" }).click();
    await expectModalVisible(modal);
    await expect(modal.getByLabel("Update interval")).toHaveValue("10");
    await expect(modal.getByLabel("Operating voltage")).toHaveValue("240");
    await expect(modal.getByLabel("Residual power")).toHaveValue("150");
  });
});
