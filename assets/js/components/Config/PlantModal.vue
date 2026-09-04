<template>
	<YamlModal
		name="plant"
		:title="$t('config.system.plant')"
		:description="$t('config.system.plantDescription')"
		docs="/reference/configuration/site#telemetry"
		defaultYaml=""
		endpoint="/config/plant"
		data-testid="plant-modal"
		@changed="$emit('changed')"
	/>
</template>

<script>
import YamlModal from "./YamlModal.vue";

export default {
	name: "PlantModal",
	components: { YamlModal },
	emits: ["changed"],
};
</script>
