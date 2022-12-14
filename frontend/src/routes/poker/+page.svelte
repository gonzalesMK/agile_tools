<script lang="ts">
	import Bees from '$lib/images/Bees.svelte';
	import Beer from '$lib/images/Beer.svelte';
	import Dash from '$lib/images/Dash.svelte';
	import QuestionMark from '$lib/images/QuestionMark.svelte';
	import { onMount } from 'svelte';
	import { Convert, type Player, type Players } from '$lib/playersDto';
	import nameStore from '$lib/stores/name';
	import playerStore from '$lib/stores/players';
	import { page } from '$app/stores';
	import { env } from '$env/dynamic/public';
	import { get } from 'svelte/store';
	const SERVER = env.PUBLIC_API_URL;
	let buttonSelected = new Array(9).fill(0);

	let gEventSource: EventSource;

	type PlayerResponse = { id: number };

	let name = $nameStore;
	let timeoutId: NodeJS.Timeout;

	let player: Player = {
		id: 0,
		name: 'juliano',
		status: -2,
		room: 0
	};
	let state: Players = {
		players: []
	};

	const buttons = ['?', 'B', '0', '1', '2', '3', '5', '8', '13'];

	function changeClickedButton(n: number) {
		buttonSelected = buttonSelected.map((val, idx) => idx == n);

		player.status = n;
		player = player;
	}

	function onClear() {
		buttonSelected = buttonSelected.fill(0);
		fetch(SERVER + 'clear', {
			method: 'POST',
			body: JSON.stringify({
				room: player.room
			}),
			headers: {
				'Content-Type': 'application/json'
			}
		});
	}

	function UpdatePlayer(p: Player) {
		if (p.id == 0 || p.room == 0) {
			return;
		}
		fetch(SERVER + 'player', {
			method: 'POST',
			body: JSON.stringify({
				id: p.id,
				room: p.room,
				status: p.status,
				name: p.name
			}),
			headers: {
				'Content-Type': 'application/json'
			}
		});
	}

	function onShow() {
		// send show to backend
		const res = fetch(SERVER + 'room', {
			method: 'POST',
			body: JSON.stringify({
				id: player.room,
				show: true
			}),
			headers: {
				'Content-Type': 'application/json'
			}
		});
	}

	function jsonIsPlayerResponseType(o: any): o is PlayerResponse {
		return 'id' in o;
	}

	function getEventSource(url: string) {
		let eventSource = new EventSource(url);

		eventSource.onmessage = ({ data }) => {
			const o: JSON = JSON.parse(data);

			if (jsonIsPlayerResponseType(o)) {
				player.id = o.id;
				playerStore.update((map) => {
					map.set(player.room.toString(), player.id.toString());
					return map;
				});
			} else if ('players' in o) {
				state = Convert.toPlayers(data);
			}
		};

		eventSource.onerror = (error) => {
			setTimeout(() => {
				gEventSource = getEventSource(url);
			}, 300);
		};

		return eventSource;
	}
	function setupSSE() {
		player.room = Number($page.url.searchParams.get('roomId'));
		const map = get(playerStore);

		let subscribe_url =
			SERVER + 'sse?room=' + encodeURIComponent(player.room) + '&name=' + encodeURIComponent(name);

		const playerId = map.get(player.room.toString()) || '';

		if (playerId != '') {
			subscribe_url += '&player=' + playerId;
		}
		let source = getEventSource(subscribe_url);

		return source;
	}
	onMount(() => {
		gEventSource = setupSSE();
	});

	function updateName(name: string) {
		clearTimeout(timeoutId);
		nameStore.set(name);
		if (name == '') {
			return;
		}

		player.name = name;
		timeoutId = setTimeout(() => {
			player = player;
		}, 2000);
	}

	$: updateName(name);
	$: UpdatePlayer(player);
</script>

<svelte:head>
	<title>Home</title>
	<meta name="description" content="Svelte demo app" />
</svelte:head>

<section>
	<input
		type="text"
		class="
	  form-control
	  block
	  px-3
	  py-1.5
	  text-base
	  font-normal
	  text-gray-700/50
	  bg-white/10 bg-clip-padding
	  border-b border-solid border-gray-300
	  rounded
	  transition
	  ease-in-out
	  m-0
	  focus:text-gray-700/80 focus:bg-white/30 focus:border-blue-600 focus:outline-none
	"
		id="exampleFormControlInput1"
		placeholder="Your name here"
		bind:value={name}
	/>
	<div class="flex flex-col overflow-hidden py-6 sm:py-12 max-w-xl m-auto">
		<div
			class="relative bg-white/60 px-6 py-5 shadow-md ring-1 ring-gray-900/5 w-full sm:rounded-lg sm:px-10"
		>
			<div class="flex flex-wrap justify-around" role="group">
				{#each buttons as btn, i}
					<button
						class="item"
						class:bg-yellow-200={buttonSelected[i]}
						class:shadow-5xl={buttonSelected[i]}
						on:click={() => changeClickedButton(i)}
					>
						{#if btn == '?'}
							<QuestionMark />
						{:else if btn == 'B'}
							<Beer />
						{:else}
							<p>{btn}</p>
						{/if}
					</button>
				{/each}
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
					{#each state.players as p}
						<tr>
							{#if p.id == player.id}
								<td>{player.name}</td>
							{:else}
								<td>{p.name}</td>
							{/if}
							<td>
								<div
									class="my-1 flex h-[3em] w-[2.5em] place-items-center justify-center rounded-xl border-2 border-solid border-y-yellow-400 border-x-yellow-500 shadow-xl ring-1 ring-gray-900/5"
								>
									{#if p.status == -2}
										<Dash />
									{:else if p.status == -1}
										<Bees />
									{:else if buttons[p.status] == 'B'}
										<Beer />
									{:else}
										<p>{buttons[p.status]}</p>
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
