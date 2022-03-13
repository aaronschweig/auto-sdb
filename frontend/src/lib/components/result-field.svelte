<script lang="ts">
	export let name = '';
	export let value = '';
	export let invalid = false;

	$: label = name[0].toUpperCase() + name.slice(1);
</script>

<label class="block">
	<span class="text-gray-700 dark:text-gray-200">{label}</span>
	<input
		{value}
		{name}
		readonly
		type="text"
		class="block w-full rounded-md bg-gray-100 border-transparent focus:border-gray-500 focus:bg-white focus:ring-0 dark:bg-gray-700"
		class:border-red-500={invalid}
	/>
</label>
