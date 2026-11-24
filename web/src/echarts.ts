// Register only the chart types Prism uses instead of shipping all of ECharts.
import * as echarts from 'echarts/core'
import { BarChart, HeatmapChart, LineChart, ScatterChart } from 'echarts/charts'
import { GridComponent, TooltipComponent, VisualMapComponent } from 'echarts/components'
import { CanvasRenderer } from 'echarts/renderers'

echarts.use([
  BarChart, HeatmapChart, LineChart, ScatterChart,
  GridComponent, TooltipComponent, VisualMapComponent, CanvasRenderer,
])

export { echarts }
export type { EChartsCoreOption } from 'echarts/core'
