<script>
	import '../app.css';
	import * as auth from '$lib/auth';
	import { onMount } from 'svelte';

	onMount(async () => {
		await auth.init()
	});
</script>

<main
	class="min-h-screen bg-gray-100 py-6 px-[3.5rem] flex flex-col justify-center align-center dark:bg-gray-800 dark:text-gray-50"
>
	<slot />
</main>
