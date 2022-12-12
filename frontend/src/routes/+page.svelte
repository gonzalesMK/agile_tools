<script lang="ts">
	import { goto } from '$app/navigation';
	import nameStore from '$lib/stores/name';

	let roomId = Math.floor(Math.random() * 1000);
	let name = $nameStore;

	$: nameStore.set(name);

	function onClick() {
		goto(roomId.toString());
	}
</script>

<svelte:head>
	<title>Home</title>
	<meta name="description" content="Svelte demo app" />
</svelte:head>

<section>
	<div
		class="relative flex min-h-screen text-gray-800 antialiased flex-col justify-center overflow-hidden py-6 sm:py-12"
	>
		<div class="relative py-3 sm:w-96 mx-auto text-center">
			<span class="text-2xl font-light ">Select a Room to Join</span>
			<div class="mt-4 bg-white shadow-md rounded-lg text-left">
				<div class="h-2 bg-yellow-400 rounded-t-md" />
				<form class="px-8 py-6 " on:submit|preventDefault={onClick}>
					<label class="block font-semibold">
						Your name
						<input
							type="text"
							placeholder="Your name here..."
							bind:value={name}
							class="border w-full h-5 px-3 py-5 mt-2 hover:outline-none focus:outline-none focus:ring-yellow-500 focus:ring-1 rounded-md"
						/>
					</label>
					<label class="block mt-3 font-semibold">
						Room id
						<input
							type="number"
							class="border w-full h-5 px-3 py-5 mt-2 hover:outline-none focus:outline-none focus:ring-yellow-500 focus:ring-1 rounded-md"
							bind:value={roomId}
						/></label
					>
					<div class="flex justify-between items-baseline">
						<button
							type="submit"
							class="mt-4 bg-yellow-400 text-white py-2 px-6 rounded-md hover:bg-yellow-500 "
							>Go</button
						>
					</div>
				</form>
			</div>
		</div>
	</div>
</section>

<style>
	section {
		display: flex;
		flex-direction: column;
		justify-content: center;
		align-items: center;
		flex: 0.6;
	}
</style>
