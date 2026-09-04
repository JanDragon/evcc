<template>
	<YamlModal
		name="levels"
		:title="$t('config.system.logLevels')"
		:description="$t('config.system.logLevelsDescription')"
		docs="/reference/configuration/log"
		:defaultYaml="defaultYaml"
		removeKey="levels"
		endpoint="/config/levels"
		data-testid="levels-modal"
		@changed="$emit('changed')"
	/>
</template>

<script>
import YamlModal from "./YamlModal.vue";
import defaultYaml from "./defaultYaml/levels.yaml?raw";

export default {
	name: "LevelsModal",
	components: { YamlModal },
	emits: ["changed"],
	data() {
		return { defaultYaml: defaultYaml.trim() };
	},
};
</script>
