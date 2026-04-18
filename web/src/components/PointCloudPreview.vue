<template>
  <div class="pointcloud-shell">
    <canvas ref="canvasRef" class="pointcloud-canvas"></canvas>
    <div v-if="statusText" class="pointcloud-status">{{ statusText }}</div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref, watch } from 'vue';

const props = defineProps<{
  src: string;
}>();

const canvasRef = ref<HTMLCanvasElement>();
const statusText = ref('正在加载点云预览...');

onMounted(() => {
  void renderPointCloud();
});

watch(
  () => props.src,
  () => {
    void renderPointCloud();
  },
);

async function renderPointCloud() {
  const canvas = canvasRef.value;
  if (!canvas || !props.src) return;

  const context = canvas.getContext('2d');
  if (!context) return;

  const width = 360;
  const height = 220;
  canvas.width = width;
  canvas.height = height;
  context.clearRect(0, 0, width, height);
  context.fillStyle = '#f5f8fc';
  context.fillRect(0, 0, width, height);

  try {
    const response = await fetch(props.src);
    const text = await response.text();
    const points = parsePointCloud(text);
    if (!points.length) {
      statusText.value = '当前点云文件不支持解析预览';
      return;
    }

    const normalized = normalizePoints(points);
    context.fillStyle = '#0f4fae';
    normalized.forEach((point) => {
      const x = point[0] * (width - 40) + 20;
      const y = point[1] * (height - 40) + 20;
      context.beginPath();
      context.arc(x, height - y, 1.8, 0, Math.PI * 2);
      context.fill();
    });
    statusText.value = `已渲染 ${normalized.length} 个点`;
  } catch {
    statusText.value = '点云预览加载失败';
  }
}

function parsePointCloud(text: string) {
  const lines = text.split(/\r?\n/).map((line) => line.trim()).filter(Boolean);
  if (!lines.length) return [] as number[][];

  if (lines[0] === 'ply') {
    return parsePly(lines);
  }
  if (lines[0].startsWith('VERSION') || lines.some((line) => line.startsWith('FIELDS'))) {
    return parsePcd(lines);
  }
  if (lines.some((line) => line.startsWith('v '))) {
    return parseObj(lines);
  }
  return [];
}

function parsePly(lines: string[]) {
  const points: number[][] = [];
  const start = lines.findIndex((line) => line === 'end_header');
  if (start < 0) return points;
  for (const line of lines.slice(start + 1)) {
    const values = line.split(/\s+/).map(Number);
    if (values.length >= 3 && values.every((value) => Number.isFinite(value))) {
      points.push(values.slice(0, 3));
    }
    if (points.length >= 1000) break;
  }
  return points;
}

function parseObj(lines: string[]) {
  const points: number[][] = [];
  for (const line of lines) {
    if (!line.startsWith('v ')) continue;
    const values = line.slice(2).split(/\s+/).map(Number);
    if (values.length >= 3 && values.every((value) => Number.isFinite(value))) {
      points.push(values.slice(0, 3));
    }
    if (points.length >= 1000) break;
  }
  return points;
}

function parsePcd(lines: string[]) {
  const points: number[][] = [];
  const start = lines.findIndex((line) => line.toUpperCase() === 'DATA ASCII');
  if (start < 0) return points;
  for (const line of lines.slice(start + 1)) {
    const values = line.split(/\s+/).map(Number);
    if (values.length >= 3 && values.every((value) => Number.isFinite(value))) {
      points.push(values.slice(0, 3));
    }
    if (points.length >= 1000) break;
  }
  return points;
}

function normalizePoints(points: number[][]) {
  const xs = points.map((point) => point[0]);
  const ys = points.map((point) => point[1]);
  const minX = Math.min(...xs);
  const maxX = Math.max(...xs);
  const minY = Math.min(...ys);
  const maxY = Math.max(...ys);
  const rangeX = maxX - minX || 1;
  const rangeY = maxY - minY || 1;

  return points.map((point) => [
    (point[0] - minX) / rangeX,
    (point[1] - minY) / rangeY,
  ]);
}
</script>

<style scoped>
.pointcloud-shell {
  display: grid;
  gap: 8px;
}

.pointcloud-canvas {
  width: 100%;
  border-radius: 12px;
  border: 1px solid var(--line);
  background: #f5f8fc;
}

.pointcloud-status {
  color: var(--text-secondary);
  font-size: 12px;
}
</style>
