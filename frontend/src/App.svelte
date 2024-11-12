<script lang="ts" module>
    let themeState: "dark" | "light" = $state("light");

    export let theme = {
        get state(): "dark" | "light" {
            return themeState;
        },
    };
</script>

<script lang="ts">
    import Moon from "./assets/icons/Moon.svelte";
    import Sun from "./assets/icons/Sun.svelte";
    import { onMount } from "svelte";
    import {
        GetUserColorTheme,
        SetUserColorTheme,
    } from "../wailsjs/go/main/App";
    import { fade } from "svelte/transition";

    const htmlElement = document.getElementsByTagName("html")[0];
    htmlElement.classList.add("light");

    function toggleColorScheme() {
        if (themeState == "dark") {
            themeState = "light";
        } else {
            themeState = "dark";
        }
        SetUserColorTheme(themeState);
    }

    onMount(() => {
        GetUserColorTheme().then((e) => (themeState = e));
    });

    $effect(() => {
        let html = document.getElementsByTagName("html")[0];
        if (themeState == "dark") {
            html.classList.replace("light", "dark");
        } else {
            html.classList.replace("dark", "light");
        }
    });
</script>

<header>
    <button class="light-dark-toggle" onclick={toggleColorScheme}>
        {#if themeState == "light"}
            <Sun></Sun>
        {:else}
            <Moon></Moon>
        {/if}
    </button>
</header>

<h1>Hello, world!</h1>

{#key theme.state}
    <!-- Transitions happen because of the key changing, not because of #if directive -->
    <p in:fade={{ delay: 200, duration: 200 }} out:fade={{ duration: 200 }}>
        {#if theme.state == "light"}
            It's day time right now!
        {:else}
            It's night time! 🎃
        {/if}
    </p>
{/key}

<style lang="postcss">
    header {
        display: flex;
        flex-flow: row;
    }

    .light-dark-toggle {
        flex-grow: 1;
        margin-left: auto;
        max-width: var(--size-8);
        font-size: var(--size-8);
        transition: var(--default-transition);
        background-color: transparent;
    }

    :root.light .light-dark-toggle {
        color: var(--stone-12);
    }

    :root.light .light-dark-toggle:hover {
        color: var(--stone-8);
    }

    :root.dark .light-dark-toggle {
        color: var(--violet-0);
    }

    :root.dark .light-dark-toggle:hover {
        color: var(--violet-2);
    }
</style>
