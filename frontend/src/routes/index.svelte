<script lang="ts">
	let loading = false;
	let isInitial = true;
	let percentage = 0.0;

	type ExtractorResult = {
		bezeichnung: string;
		lagerklasse: string;
		signalwort: string;
		hSaezte: string[];
		pSaezte: string[];
		ghs: string;
		wgk: string;
	};

	let result = {
		bezeichnung: '',
		lagerklasse: '',
		signalwort: '',
		ghs: '',
		wgk: '',
		hSaezte: [],
		pSaezte: []
	} as ExtractorResult;

	async function extract(e: SubmitEvent) {
		const form = new FormData(e.target as HTMLFormElement);

		loading = true;

		const extractorResult = await fetch('/extract', {
			method: 'POST',
			body: form
		}).then((res) => res.json() as Promise<ExtractorResult>);

		result = extractorResult;

		loading = false;
		isInitial = false;

		const total = Object.keys(extractorResult).length;
		let notExtracted = 0;

		Object.keys(extractorResult).forEach((key) => {
			if (extractorResult[key]?.length == 0) {
				notExtracted += 1;
			}
		});

		percentage = parseFloat((((total - notExtracted) / total) * 100).toFixed(2));
	}
</script>

<h1 class="text-4xl py-4">SDB-Extractor</h1>
<div class="grid grid-cols-2 gap-4">
	<form
		class="flex flex-col justify-center items-center gap-y-4 px-4 bg-white dark:bg-gray-500 shadow-lg rounded-3xl"
		on:submit|preventDefault={extract}
	>
		<input
			type="file"
			name="file"
			class="block w-full text-sm text-gray-900 bg-gray-50 rounded-lg border border-gray-300 cursor-pointer dark:text-gray-400 focus:outline-none focus:border-transparent dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400"
		/>
		<button
			class="px-3 py-2 bg-blue-600 hover:bg-blue-400 dark:bg-sky-600 dark:hover:bg-sky-400 transition-colors rounded shadow text-white text-lg w-72"
			type="submit"
		>
			Extract
		</button>
	</form>
	<div class="bg-white dark:bg-gray-500 flex justify-center px-4 py-10 rounded-3xl shadow-lg">
		<svg
			class="animate-spin h-16 text-green-700 dark:text-teal-200"
			class:hidden={!loading}
			xmlns="http://www.w3.org/2000/svg"
			fill="none"
			viewBox="0 0 24 24"
		>
			<circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
			<path
				class="opacity-75"
				fill="currentColor"
				d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
			/>
		</svg>
		<div class="w-full space-y-2 dark:text-gray-200" class:hidden={loading}>
			<div
				class="w-full bg-gray-200 dark:bg-gray-200 rounded overflow-hidden"
				class:hidden={isInitial}
			>
				<div
					class="bg-green-600 text-xs font-medium text-center h-2 leading-none"
					style="width: {percentage}%"
				/>
			</div>
			<label class="block">
				<span class="text-gray-700 dark:text-gray-200">Bezeichnung</span>
				<input
					value={result.bezeichnung}
					name="bezeichnung"
					readonly
					type="text"
					class="block w-full rounded-md bg-gray-100 border-transparent focus:border-gray-500 focus:bg-white focus:ring-0 dark:bg-gray-700"
					class:border-red-500={!isInitial && !result.bezeichnung?.length}
				/>
			</label>
			<label class="block">
				<span class="text-gray-700 dark:text-gray-200">Lagerklasse</span>
				<input
					value={result.lagerklasse}
					name="lagerklasse"
					readonly
					type="text"
					class="block w-full rounded-md bg-gray-100 border-transparent focus:border-gray-500 focus:bg-white focus:ring-0 dark:bg-gray-700"
					class:border-red-500={!isInitial && !result.lagerklasse?.length}
				/>
			</label>
			<label class="block">
				<span class="text-gray-700 dark:text-gray-200">Signalwort</span>
				<input
					value={result.signalwort}
					name="signalwort"
					readonly
					type="text"
					class="block w-full rounded-md bg-gray-100 border-transparent focus:border-gray-500 focus:bg-white focus:ring-0 dark:bg-gray-700"
					class:border-red-500={!isInitial && !result.signalwort?.length}
				/>
			</label>
			<label class="block">
				<span class="text-gray-700 dark:text-gray-200">GHS</span>
				<input
					value={result.ghs}
					name="ghs"
					readonly
					type="text"
					class="block w-full rounded-md bg-gray-100 border-transparent focus:border-gray-500 focus:bg-white focus:ring-0 dark:bg-gray-700"
					class:border-red-500={!isInitial && !result.ghs?.length}
				/>
			</label>
			<label class="block">
				<span class="text-gray-700 dark:text-gray-200">H-Sätze</span>
				<input
					value={result.hSaezte}
					name="hSaezte"
					readonly
					type="text"
					class="block w-full rounded-md bg-gray-100 border-transparent focus:border-gray-500 focus:bg-white focus:ring-0 dark:bg-gray-700"
					class:border-red-500={!isInitial && !result.hSaezte?.length}
				/>
			</label>
			<label class="block">
				<span class="text-gray-700 dark:text-gray-200">P-Sätze</span>
				<input
					value={result.pSaezte}
					name="pSaezte"
					readonly
					type="text"
					class="block w-full rounded-md bg-gray-100 border-transparent focus:border-gray-500 focus:bg-white focus:ring-0 dark:bg-gray-700"
					class:border-red-500={!isInitial && !result.pSaezte?.length}
				/>
			</label>
			<label class="block">
				<span class="text-gray-700 dark:text-gray-200">WGK</span>
				<input
					value={result.wgk}
					name="wgk"
					readonly
					type="text"
					class="block w-full rounded-md bg-gray-100 border-transparent focus:border-gray-500 focus:bg-white focus:ring-0 dark:bg-gray-700"
					class:border-red-500={!isInitial && !result.wgk?.length}
				/>
			</label>
		</div>
	</div>
</div>
