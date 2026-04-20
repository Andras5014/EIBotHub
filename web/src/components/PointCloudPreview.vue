<template>
  <div class="pointcloud-shell">
    <div class="viewer-toolbar">
      <a-space wrap size="small">
        <a-button size="small" @click="resetView">重置视角</a-button>
        <a-button size="small" @click="fitView">适配视图</a-button>
      </a-space>
      <span class="viewer-hint">左键拖拽旋转，Shift + 拖拽平移，滚轮缩放</span>
    </div>

    <div class="canvas-wrap">
      <canvas
        ref="canvasRef"
        class="pointcloud-canvas"
        @pointerdown="handlePointerDown"
        @pointermove="handlePointerMove"
        @pointerup="handlePointerUp"
        @pointerleave="handlePointerUp"
        @wheel.prevent="handleWheel"
        @contextmenu.prevent
      ></canvas>
      <div v-if="overlayText" class="pointcloud-overlay">{{ overlayText }}</div>
    </div>

    <div class="pointcloud-status">
      <span>{{ statusText }}</span>
      <span v-if="sampled">{{ sampleText }}</span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, nextTick, onBeforeUnmount, onMounted, ref, watch } from 'vue';

type Point3D = {
  x: number;
  y: number;
  z: number;
};

type ProjectedPoint = {
  x: number;
  y: number;
  depth: number;
  size: number;
  color: string;
};

type ViewState = {
  yaw: number;
  pitch: number;
  zoom: number;
  panX: number;
  panY: number;
};

const props = defineProps<{
  src: string;
}>();

const DEFAULT_WIDTH = 420;
const DEFAULT_HEIGHT = 280;
const MAX_RENDER_POINTS = 12000;
const canvasRef = ref<HTMLCanvasElement>();
const statusText = ref('正在加载点云预览...');
const overlayText = ref('加载中...');
const sampled = ref(false);
const pointCount = ref(0);
const viewState = ref<ViewState>(createDefaultView());

let points: Point3D[] = [];
let normalizedPoints: Point3D[] = [];
let resizeObserver: ResizeObserver | undefined;
let dragging = false;
let dragMode: 'rotate' | 'pan' = 'rotate';
let lastPointer = { x: 0, y: 0 };

const sampleText = computed(() => {
  if (!sampled.value) {
    return '';
  }
  return `已采样渲染 ${pointCount.value} 个点`;
});

onMounted(async () => {
  await nextTick();
  bindResizeObserver();
  await renderPointCloud();
});

onBeforeUnmount(() => {
  resizeObserver?.disconnect();
});

watch(
  () => props.src,
  async () => {
    await renderPointCloud();
  },
);

function createDefaultView(): ViewState {
  return {
    yaw: -0.55,
    pitch: 0.4,
    zoom: 1,
    panX: 0,
    panY: 0,
  };
}

function resetView() {
  viewState.value = createDefaultView();
  drawScene();
}

function fitView() {
  viewState.value.panX = 0;
  viewState.value.panY = 0;
  viewState.value.zoom = 1;
  drawScene();
}

function bindResizeObserver() {
  if (!canvasRef.value || typeof ResizeObserver === 'undefined') {
    return;
  }
  resizeObserver = new ResizeObserver(() => {
    syncCanvasSize();
    drawScene();
  });
  resizeObserver.observe(canvasRef.value);
}

function syncCanvasSize() {
  const canvas = canvasRef.value;
  if (!canvas) return;
  const rect = canvas.getBoundingClientRect();
  canvas.width = Math.max(Math.floor(rect.width || DEFAULT_WIDTH), DEFAULT_WIDTH);
  canvas.height = Math.max(Math.floor(rect.height || DEFAULT_HEIGHT), DEFAULT_HEIGHT);
}

async function renderPointCloud() {
  const canvas = canvasRef.value;
  if (!canvas || !props.src) return;

  overlayText.value = '加载中...';
  statusText.value = '正在加载点云预览...';
  sampled.value = false;
  pointCount.value = 0;
  points = [];
  normalizedPoints = [];
  resetView();
  syncCanvasSize();
  drawPlaceholder();

  try {
    const response = await fetch(props.src);
    const text = await response.text();
    points = parsePointCloud(text);
    if (!points.length) {
      overlayText.value = '当前点云文件不支持解析预览';
      statusText.value = '无法识别点云格式，请直接下载原始文件';
      return;
    }

    const sampledPoints = samplePoints(points, MAX_RENDER_POINTS);
    sampled.value = sampledPoints.length < points.length;
    pointCount.value = sampledPoints.length;
    normalizedPoints = normalizePoints(sampledPoints);
    overlayText.value = '';
    statusText.value = `3D 预览已就绪，共解析 ${points.length} 个点`;
    drawScene();
  } catch {
    overlayText.value = '点云预览加载失败';
    statusText.value = '点云预览加载失败，请稍后重试或直接下载原始文件';
  }
}

function drawPlaceholder() {
  const canvas = canvasRef.value;
  if (!canvas) return;
  const context = canvas.getContext('2d');
  if (!context) return;
  context.clearRect(0, 0, canvas.width, canvas.height);
  context.fillStyle = '#f5f8fc';
  context.fillRect(0, 0, canvas.width, canvas.height);
}

function drawScene() {
  const canvas = canvasRef.value;
  if (!canvas) return;
  const context = canvas.getContext('2d');
  if (!context) return;

  syncCanvasSize();
  context.clearRect(0, 0, canvas.width, canvas.height);
  drawBackground(context, canvas.width, canvas.height);

  if (!normalizedPoints.length) {
    return;
  }

  const projected = normalizedPoints
    .map((point) => projectPoint(point, canvas.width, canvas.height))
    .filter((point): point is ProjectedPoint => point !== null)
    .sort((left, right) => left.depth - right.depth);

  for (const point of projected) {
    context.beginPath();
    context.fillStyle = point.color;
    context.arc(point.x, point.y, point.size, 0, Math.PI * 2);
    context.fill();
  }
}

function drawBackground(context: CanvasRenderingContext2D, width: number, height: number) {
  const gradient = context.createLinearGradient(0, 0, width, height);
  gradient.addColorStop(0, '#f8fbff');
  gradient.addColorStop(1, '#eef4fb');
  context.fillStyle = gradient;
  context.fillRect(0, 0, width, height);

  context.strokeStyle = 'rgba(22, 50, 79, 0.08)';
  context.lineWidth = 1;
  for (let index = 1; index < 6; index++) {
    const y = (height / 6) * index;
    context.beginPath();
    context.moveTo(0, y);
    context.lineTo(width, y);
    context.stroke();
  }
}

function handlePointerDown(event: PointerEvent) {
  dragging = true;
  dragMode = event.shiftKey || event.button === 1 || event.button === 2 ? 'pan' : 'rotate';
  lastPointer = { x: event.clientX, y: event.clientY };
}

function handlePointerMove(event: PointerEvent) {
  if (!dragging) return;

  const deltaX = event.clientX - lastPointer.x;
  const deltaY = event.clientY - lastPointer.y;
  lastPointer = { x: event.clientX, y: event.clientY };

  if (dragMode === 'rotate') {
    viewState.value.yaw += deltaX * 0.01;
    viewState.value.pitch = clamp(viewState.value.pitch + deltaY * 0.01, -1.4, 1.4);
  } else {
    viewState.value.panX += deltaX;
    viewState.value.panY += deltaY;
  }
  drawScene();
}

function handlePointerUp() {
  dragging = false;
}

function handleWheel(event: WheelEvent) {
  const nextZoom = event.deltaY > 0 ? viewState.value.zoom * 0.92 : viewState.value.zoom * 1.08;
  viewState.value.zoom = clamp(nextZoom, 0.3, 4);
  drawScene();
}

function projectPoint(point: Point3D, width: number, height: number): ProjectedPoint | null {
  const rotated = rotatePoint(point, viewState.value.yaw, viewState.value.pitch);
  const cameraDistance = 3.2 / viewState.value.zoom;
  const perspective = cameraDistance / (cameraDistance + rotated.z + 1.8);
  const x = rotated.x * perspective * (width * 0.32) + width / 2 + viewState.value.panX;
  const y = rotated.y * perspective * (height * 0.32) + height / 2 + viewState.value.panY;

  if (!Number.isFinite(x) || !Number.isFinite(y)) {
    return null;
  }

  const depth = rotated.z;
  const size = Math.max(1.1, perspective * 2.4);
  const hue = 205 + Math.round((rotated.z + 1) * 32);
  return {
    x,
    y,
    depth,
    size,
    color: `hsla(${hue}, 78%, 42%, 0.85)`,
  };
}

function rotatePoint(point: Point3D, yaw: number, pitch: number): Point3D {
  const cosYaw = Math.cos(yaw);
  const sinYaw = Math.sin(yaw);
  const cosPitch = Math.cos(pitch);
  const sinPitch = Math.sin(pitch);

  const x1 = point.x * cosYaw - point.z * sinYaw;
  const z1 = point.x * sinYaw + point.z * cosYaw;
  const y2 = point.y * cosPitch - z1 * sinPitch;
  const z2 = point.y * sinPitch + z1 * cosPitch;

  return { x: x1, y: y2, z: z2 };
}

function parsePointCloud(text: string) {
  const lines = text.split(/\r?\n/).map((line) => line.trim()).filter(Boolean);
  if (!lines.length) return [] as Point3D[];

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
  const points: Point3D[] = [];
  const start = lines.findIndex((line) => line === 'end_header');
  if (start < 0) return points;
  for (const line of lines.slice(start + 1)) {
    const values = line.split(/\s+/).map(Number);
    if (values.length >= 3 && values.slice(0, 3).every((value) => Number.isFinite(value))) {
      points.push({ x: values[0], y: values[1], z: values[2] });
    }
  }
  return points;
}

function parseObj(lines: string[]) {
  const points: Point3D[] = [];
  for (const line of lines) {
    if (!line.startsWith('v ')) continue;
    const values = line.slice(2).split(/\s+/).map(Number);
    if (values.length >= 3 && values.slice(0, 3).every((value) => Number.isFinite(value))) {
      points.push({ x: values[0], y: values[1], z: values[2] });
    }
  }
  return points;
}

function parsePcd(lines: string[]) {
  const points: Point3D[] = [];
  const start = lines.findIndex((line) => line.toUpperCase() === 'DATA ASCII');
  if (start < 0) return points;
  for (const line of lines.slice(start + 1)) {
    const values = line.split(/\s+/).map(Number);
    if (values.length >= 3 && values.slice(0, 3).every((value) => Number.isFinite(value))) {
      points.push({ x: values[0], y: values[1], z: values[2] });
    }
  }
  return points;
}

function samplePoints(items: Point3D[], maxCount: number) {
  if (items.length <= maxCount) {
    return items;
  }
  const step = Math.ceil(items.length / maxCount);
  const sampledPoints: Point3D[] = [];
  for (let index = 0; index < items.length; index += step) {
    sampledPoints.push(items[index]);
  }
  return sampledPoints;
}

function normalizePoints(items: Point3D[]) {
  const xs = items.map((point) => point.x);
  const ys = items.map((point) => point.y);
  const zs = items.map((point) => point.z);
  const minX = Math.min(...xs);
  const maxX = Math.max(...xs);
  const minY = Math.min(...ys);
  const maxY = Math.max(...ys);
  const minZ = Math.min(...zs);
  const maxZ = Math.max(...zs);
  const centerX = (minX + maxX) / 2;
  const centerY = (minY + maxY) / 2;
  const centerZ = (minZ + maxZ) / 2;
  const scale = Math.max(maxX - minX, maxY - minY, maxZ - minZ) || 1;

  return items.map((point) => ({
    x: ((point.x - centerX) / scale) * 2,
    y: ((point.y - centerY) / scale) * 2,
    z: ((point.z - centerZ) / scale) * 2,
  }));
}

function clamp(value: number, min: number, max: number) {
  return Math.min(Math.max(value, min), max);
}
</script>

<style scoped>
.pointcloud-shell {
  display: grid;
  gap: 10px;
}

.viewer-toolbar {
  display: flex;
  justify-content: space-between;
  gap: 12px;
  align-items: center;
  flex-wrap: wrap;
}

.viewer-hint {
  color: var(--text-secondary);
  font-size: 12px;
}

.canvas-wrap {
  position: relative;
}

.pointcloud-canvas {
  width: 100%;
  min-height: 280px;
  border-radius: 12px;
  border: 1px solid var(--line);
  background: #f5f8fc;
  touch-action: none;
  cursor: grab;
}

.pointcloud-canvas:active {
  cursor: grabbing;
}

.pointcloud-overlay {
  position: absolute;
  inset: 0;
  display: grid;
  place-items: center;
  color: var(--text-secondary);
  font-size: 13px;
  background: rgba(248, 251, 255, 0.68);
  border-radius: 12px;
}

.pointcloud-status {
  display: flex;
  justify-content: space-between;
  gap: 12px;
  flex-wrap: wrap;
  color: var(--text-secondary);
  font-size: 12px;
}

@media (max-width: 768px) {
  .pointcloud-canvas {
    min-height: 240px;
  }
}
</style>
