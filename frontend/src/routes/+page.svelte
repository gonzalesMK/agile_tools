<script lang="ts">
	import Bees from '$lib/images/bees.svg';
	import Beer from '$lib/images/beer.svg';
	import Dash from '$lib/images/dash.svg';
	import Question from '$lib/images/question-mark.svg';
	import { onMount } from 'svelte';

	let buttonSelected = new Array(9).fill(0);

	let state = {
		players: []
	};
	// let state = {
	// 	players: [
	// 		{
	// 			name: 'Juliano Decico Negri',
	// 			status: -1
	// 		},
	// 		{
	// 			name: 'Witchy Woman',
	// 			status: -1
	// 		},
	// 		{
	// 			name: 'Shining Star',
	// 			status: -3
	// 		},
	// 		{
	// 			name: 'Shining Star',
	// 			status: 0
	// 		},
	// 		{
	// 			name: 'Shining Star',
	// 			status: 13
	// 		}
	// 	]
	// };

	function changeClickedButton(n: number) {
		buttonSelected = buttonSelected.map((val, idx) => idx == n);
		// send result to backend
	}

	function onClear() {
		buttonSelected = buttonSelected.fill(0);
		// send clear to backend
	}

	function onShow() {
		// send show to backend
	}
	const SSE_LOCAL_URL = 'http://127.0.0.1:3000/sse?room=101&name=juliano';
	onMount(() => {
		const source = new EventSource(SSE_LOCAL_URL);
		source.onopen = () => {
			console.log('event opened');
		};
		source.onmessage = ({ data }) => {
			state = JSON.parse(data);
			console.log('onmessage ', data);
		};
		source.onerror = (error) => console.log('source error ', error);

		return () => {
			if (source.readyState === 1) {
				source.close();
			}
			console.log('HERE');
		};
	});
</script>

<svelte:head>
	<title>Home</title>
	<meta name="description" content="Svelte demo app" />
</svelte:head>

<section>
	<div class="flex min-h-screen flex-col overflow-hidden py-6 sm:py-12 max-w-xl m-auto">
		<div class="relative px-6 py-5 shadow-md ring-1 ring-gray-900/5 w-full sm:rounded-lg sm:px-10">
			<div class="flex flex-wrap justify-around">
				<button
					class="item"
					class:bg-yellow-200={buttonSelected[0]}
					class:shadow-5xl={buttonSelected[0]}
					on:click={() => changeClickedButton(0)}
				>
					<img src={Question} class="w-[1.1em]" alt="Welcome" />
				</button>
				<button
					class="item"
					class:bg-yellow-200={buttonSelected[1]}
					class:shadow-5xl={buttonSelected[1]}
					on:click={() => changeClickedButton(1)}
				>
					<img src={Beer} alt="Welcome" />
				</button>
				<button
					class="item"
					class:bg-yellow-200={buttonSelected[2]}
					class:shadow-5xl={buttonSelected[2]}
					on:click={() => changeClickedButton(2)}
				>
					<p>0</p>
				</button>
				<button
					class="item"
					class:bg-yellow-200={buttonSelected[3]}
					class:shadow-5xl={buttonSelected[3]}
					on:click={() => changeClickedButton(3)}
				>
					<p>1</p>
				</button>
				<button
					class="item"
					class:bg-yellow-200={buttonSelected[4]}
					class:shadow-5xl={buttonSelected[4]}
					on:click={() => changeClickedButton(4)}
				>
					<p>2</p>
				</button>
				<button
					class="item"
					class:bg-yellow-200={buttonSelected[5]}
					class:shadow-5xl={buttonSelected[5]}
					on:click={() => changeClickedButton(5)}
				>
					<p>3</p>
				</button>
				<button
					class="item"
					class:bg-yellow-200={buttonSelected[6]}
					class:shadow-5xl={buttonSelected[6]}
					on:click={() => changeClickedButton(6)}
				>
					<p>5</p>
				</button>
				<button
					class="item"
					class:bg-yellow-200={buttonSelected[7]}
					class:shadow-5xl={buttonSelected[7]}
					on:click={() => changeClickedButton(7)}
				>
					<p>8</p>
				</button>
				<button
					class="item"
					class:bg-yellow-200={buttonSelected[8]}
					class:shadow-5xl={buttonSelected[8]}
					on:click={() => changeClickedButton(8)}
				>
					<p>13</p>
				</button>
			</div>
		</div>

		<div class="flex mt-10 justify-center space-x-2 w-full sm:rounded-lg sm:px-10">
			<button
				type="button"
				class="inline-block rounded bg-blue-600 px-6 py-2.5 text-xs font-medium uppercase leading-tight text-white shadow-md transition duration-150 ease-in-out hover:bg-blue-700 hover:shadow-lg focus:bg-blue-700 focus:shadow-lg focus:outline-none focus:ring-0 active:bg-blue-800 active:shadow-lg"
				on:click={onShow}>Show Results</button
			>
			<button
				type="button"
				class="inline-block rounded bg-blue-600 px-6 py-2.5 text-xs font-medium uppercase leading-tight text-white shadow-md transition duration-150 ease-in-out hover:bg-blue-700 hover:shadow-lg focus:bg-blue-700 focus:shadow-lg focus:outline-none focus:ring-0 active:bg-blue-800 active:shadow-lg"
				on:click={onClear}>Clear Results</button
			>
		</div>

		<div
			class="mt-10 bg-white/60 px-6 pt-10 shadow-xl ring-1 ring-gray-900/5 sm:mx-auto w-full m:rounded-lg sm:px-10"
		>
			<p class="pb-10 text-center text-2xl font-bold">Results</p>
			<table class="min-w-full divide-y divide-gray-300/50">
				<thead class="text-left">
					<tr>
						<th>Name</th>
						<th>Story Points</th>
					</tr>
				</thead>
				<tbody class="divide-y divide-gray-300/50 leading-10">
					{#each state.players as player}
						<tr>
							<td>{player.name}</td>
							<td>
								<div
									class="my-1 flex h-[3em] w-[2.5em] place-items-center justify-center rounded-xl border-2 border-solid border-y-yellow-400 border-x-yellow-500 shadow-xl ring-1 ring-gray-900/5"
								>
									{#if player.status == -1}
										<img src={Beer} alt="Welcome" />
									{:else if player.status == -2}
										<img src={Dash} alt="Welcome" />
									{:else if player.status == -3}
										<img src={Bees} alt="Welcome" class="p-1 pr-2" />
									{:else}
										<p>{player.status}</p>
									{/if}
								</div>
							</td>
						</tr>
					{/each}
				</tbody>
			</table>
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
