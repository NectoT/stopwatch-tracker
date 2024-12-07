<script lang="ts">
    import { GetStopwatches } from "@wails/go_timer_tracker/app";
    import { ColorTheme, type StopwatchData } from "@wails/go_timer_tracker/models";
    import { onMount } from "svelte";
    import { Chart, type ChartConfiguration } from "chart.js/auto";
    import type { ActionReturn } from "svelte/action";
    import { theme, themeColors } from "@App";
    
    type TimeInterval = "Day" | "Week" | "Month";
    type TimeIntervalData = {
        [timeInterval in TimeInterval]: {
            timestep: number, 
            startDate: Date, 
            timestepLabels: () => string[]
        }
    };
    let selectedTimePeriod: TimeInterval = $state("Day");
    let timeIntervalData: TimeIntervalData = {
        "Day": {
            timestep: 1000 * 60 * 60,
            startDate: new Date(Date.now() - 1000 * 60 * 60 * 24),
            timestepLabels: () => {
                let labels: string[] = []
                let tpd = timeIntervalData[selectedTimePeriod];
                for (let i = 0; i < 24; i++) {
                    labels.push(new Date(tpd.startDate.getTime() + i * tpd.timestep).getHours().toString())
                }

                return labels;
            }
        },
        "Week": {
            timestep: 1000 * 60 * 60 * 24,
            startDate: new Date(Date.now() - 1000 * 60 * 60 * 24 * 7),
            timestepLabels: () => {
                let labels: string[] = []
                let tpd = timeIntervalData[selectedTimePeriod];
                for (let i = 0; i < 7; i++) {
                    switch (new Date(tpd.startDate.getTime() + i * tpd.timestep).getDay()) {
                        case 0:
                            labels.push('Mon');
                            break;
                        case 1:
                            labels.push('Tue');
                            break;
                        case 2:
                            labels.push('Wed');
                            break;
                        case 3:
                            labels.push('Thu');
                            break;
                        case 4:
                            labels.push('Fri');
                            break;
                        case 5:
                            labels.push('Sat');
                            break;
                        case 6:
                            labels.push('Sun');
                            break;
                    }
                }

                return labels;
            }
        },
        "Month": {
            timestep: 1000 * 60 * 60 * 24,
            startDate: new Date(Date.now() - 1000 * 60 * 60 * 24 * 30),
            timestepLabels: () => {
                let labels: string[] = []
                let tpd = timeIntervalData[selectedTimePeriod];
                for (let i = 0; i < 30; i++) {
                    labels.push(new Date(tpd.startDate.getTime() + i * tpd.timestep).getDate().toString())
                }

                return labels;
            }
        }
    }

    let data: StopwatchData[] = $state([]);
    
    let fontColor = $derived(themeColors.font.primary);

    let gridlineColor = $derived(theme.state === ColorTheme.Light ? '#EEEEEE' : '#444444');

    let statsEndDate = new Date(Date.now());
    let statsStartDate = $derived(timeIntervalData[selectedTimePeriod].startDate);
    let timestep = $derived(timeIntervalData[selectedTimePeriod].timestep);
    let timestepLabels = $derived(timeIntervalData[selectedTimePeriod].timestepLabels())
    
    let pieConfig: ChartConfiguration<"pie"> = $derived({
        type: 'pie',
        data: {
            datasets: [{
                data: data.map(e => getActiveRatio(
                    e.ActivePeriods ?? [], statsStartDate, statsEndDate, e.IsActive
                )),
                backgroundColor: data.map(e => `hsl(${e.Hue}, ${theme.state === "light" ? "50%" : "40%"}, 50%)`)
            }],
            labels: data.map(e => e.Name)
        },
        options: {
            rotation: 180,
            cutout: '50%',
            plugins: {
                legend: {
                    labels: {
                        color: fontColor
                    }
                }
            },
        }
    });

    let barConfig: ChartConfiguration<"bar"> = $derived({
        type: "bar",
        data: {
            labels: timestepLabels,
            datasets: data.map(d => ({
                data: groupPeriodsByFixedIntervals(
                    d.ActivePeriods ?? [], statsStartDate, statsEndDate, timestep, d.IsActive
                ),
                label: d.Name,
                stack: '1',
                borderWidth: 1,
                backgroundColor: `hsl(${d.Hue}, ${theme.state === "light" ? "50%" : "40%"}, 50%)`,
            }))
        },
        options: {
            scales: {
                x: {
                    grid: {
                        color: gridlineColor
                    }
                },
                y: {
                    grid: {
                        color: gridlineColor
                    }
                }
            },
            plugins: {
                legend: {
                    labels: {
                        color: fontColor
                    }
                }
            }
        }
    });

    function getActiveRatio(
        periods: string[][], start: Date, end: Date, currentlyActive: boolean
    ) {
        let result = groupPeriodsByFixedIntervals(
            periods, start, end, end.getTime() - start.getTime(), currentlyActive
        );
        return result.length !== 0 ? result[0] : 0;
    }

    function groupPeriodsByFixedIntervals(
        periods: string[][], start: Date, end: Date, step: number, currentlyActive: boolean
    ): number[] {
        //        s|   |   |   |   |   |   |   |   |   |   |e
        // [---]---[--------]--[--]-[--]-[---]--[-------]---> - Timeline
        // 
        // [] - periods, | | - fixed intervals, defined by `step`

        if (currentlyActive) {
            periods = periods.slice();
            periods[periods.length - 1] = [
                periods[periods.length - 1][0],
                new Date(Date.now()).toISOString()
            ];
        }

        let activityRatio: number[] = [];
        let currPeriodIdx = 0;
        for (let i = start.getTime(); i < end.getTime(); i += step) {
            // The fixed interval doesn't intersect current period and 
            // it's to the right of the period.
            while (
                currPeriodIdx !== periods.length && 
                Date.parse(periods[currPeriodIdx][0]) < i
            ) {
                currPeriodIdx += 1;
            }

            let ratio = 0;
            while (
                currPeriodIdx !== periods.length && 
                Date.parse(periods[currPeriodIdx][0]) < i + step
            ) {
                let currPeriod = periods[currPeriodIdx];
                let intersection = [
                    Math.max(i, Date.parse(currPeriod[0])), 
                    Math.min(i + step, Date.parse(currPeriod[1]))
                ];
                ratio += (intersection[1] - intersection[0]) / step;

                currPeriodIdx += 1;
            }

            // For cases where a period starts inside the fixed interval but ends outside it,
            // we shouldn't immediately go to the next period, and instead check this same
            // period in the next fixed interval
            if (currPeriodIdx !== 0 && Date.parse(periods[currPeriodIdx - 1][1]) > i + step) {
                currPeriodIdx -= 1;
            }

            activityRatio.push(ratio);
        }

        return activityRatio.map(e => e * step / 1000 / 60 / 60);
    }

    function createChart(
        node: HTMLCanvasElement, 
        configuration: ChartConfiguration<"pie" | "bar">
    ): ActionReturn<ChartConfiguration> {
        const chart = new Chart(node, configuration);
        return {
            update: (configuration: ChartConfiguration) => {
                chart.data = (configuration as ChartConfiguration<"pie" | "bar">).data;
                chart.legend!.options.labels.color = fontColor;
                if (chart.options.scales?.x !== undefined) {
                    chart.options.scales.x.grid!.color = gridlineColor;
                }
                if (chart.options.scales?.y !== undefined) {
                    chart.options.scales.y.grid!.color = gridlineColor;
                }
                chart.update();
            }
        };
    }

    onMount(async () => {
        data = Object.values(await GetStopwatches() ?? {});
    })
</script>

<div class="container">

    <div class="interval-selector">
        {#each Object.keys(timeIntervalData) as period}
            <button class="filled" class:selected={period === selectedTimePeriod} onclick={_ => selectedTimePeriod = period as TimeInterval}>
                {period}
            </button>
        {/each}
    </div>

    <div class="charts">
        <div class="pie-holder">
            <canvas class="pie" use:createChart={pieConfig}></canvas>
        </div>
        <div class="bar-holder">
            <canvas class="bar" use:createChart={barConfig}></canvas>
        </div>
    </div>
</div>

<style>
    canvas {
        width: 100%;
        height: 100%;
    }

    .container {
        max-width: 80%;
        height: 100%;
        margin: auto;
        display: flex;
        flex-flow: column;
    }

    .interval-selector {
        min-width: var(--size-14);
        width: 25%;
        height: var(--size-10);
        display: flex;
        flex-direction: row;
        column-gap: var(--size-5);
        align-self: center;
        padding: var(--size-3);
        border-radius: var(--size-1) var(--size-1) var(--size-8) var(--size-8);
        /* background-color: var(--sec-bg-color); */
        background-image: linear-gradient(var(--pr-bg-color), var(--sec-bg-color));
        margin-bottom: var(--size-4);
        filter: drop-shadow(0px 4px 8px var(--pr-font-color));
    }

    .interval-selector button {
        flex-grow: 1;
        background-color: transparent;
        font-size: large;
        font-weight: bold;
        transition: var(--default-transition);
    }
    
    .interval-selector button.selected {
        font-size: var(--size-5);
        background-image: radial-gradient(var(--pr-bg-color), transparent 80%);
    }

    .charts {
        display: grid;
        flex-grow: 1;
        width: 100%;
        grid-template-columns: 40% 60%;
        grid-template-rows: 50% 50%;
    }

    .charts>* {
        grid-row: auto / span 1;
    }

    .pie-holder {
        grid-row: 1 / span 2;
    }

    .bar-holder {
        grid-row: 1 / span 2;
        /* align-self: center; */
    }
</style>