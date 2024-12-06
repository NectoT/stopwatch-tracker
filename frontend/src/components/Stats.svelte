<script lang="ts">
    import { GetStopwatches } from "@wails/go_timer_tracker/app";
    import type { StopwatchData } from "@wails/go_timer_tracker/models";
    import { onMount } from "svelte";
    import { Chart, type ChartConfiguration } from "chart.js/auto";
    import type { ActionReturn } from "svelte/action";
    import { theme } from "@App";
    
    let data: StopwatchData[] = $state([]);

    let fontColor = $derived.by(() => {
        theme.state
        return getComputedStyle(document.body).getPropertyValue("--pr-font-color");
    });
    
    let pieConfig: ChartConfiguration<"pie"> = $derived({
        type: 'pie',
        data: {
            datasets: [{
                data: data.map(e => e.TimeAccumulated),
                backgroundColor: data.map(e => `hsl(${e.Hue}, ${theme.state === "light" ? "50%" : "40%"}, 50%)`)
            }],
            labels: data.map(e => e.Name)
        },
        options: {
            plugins: {
                legend: {
                    labels: {
                        color: fontColor
                    }
                }
            }
        }
    });

    function createPieChart(
        node: HTMLCanvasElement, 
        configuration: ChartConfiguration<"pie">
    ): ActionReturn<ChartConfiguration<"pie">> {
        const chart = new Chart(node, configuration);
        return {
            update: (configuration) => {
                chart.data = configuration.data;
                console.log(fontColor);
                chart.legend!.options.labels.color = fontColor;
                chart.update();
            }
        };
    }

    onMount(async () => {
        data = Object.values(await GetStopwatches());
    })
</script>

<div class="container">
    <canvas class="pie" use:createPieChart={pieConfig}></canvas>
</div>

<style>
    .container {
        max-width: 80%;
        margin: auto;
    }

    .pie {
        height: 100%;
        width: 100%;
    }
</style>