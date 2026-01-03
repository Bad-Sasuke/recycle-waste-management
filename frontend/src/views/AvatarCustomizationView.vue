<template>
  <div class="relative w-full h-[calc(100vh-64px)] overflow-hidden">
    <!-- 3D Scene Container -->
    <div ref="container" class="absolute inset-0"></div>

    <!-- Stats Overlay (Floating near avatar) -->
    <div
      class="absolute bottom-6 left-6 p-4 bg-white/70 backdrop-blur-md rounded-xl border border-white/40 shadow-lg w-64">
      <h3 class="text-sm font-semibold text-gray-700 mb-2 flex items-center gap-2">
        <span class="text-lg">{{ currentOutfitData?.icon }}</span>
        {{ currentOutfitData?.name || 'Select Outfit' }}
      </h3>
      <div class="space-y-2">
        <div v-for="(stat, index) in currentStats" :key="index" class="space-y-1">
          <div class="flex justify-between text-xs font-medium text-gray-600">
            <span>{{ stat.label }}</span>
            <span>{{ stat.value }}%</span>
          </div>
          <div class="h-2 bg-gray-200 rounded-full overflow-hidden">
            <div class="h-full bg-gradient-to-r from-green-400 to-emerald-500 transition-all duration-500 ease-out"
              :style="{ width: `${stat.value}%` }"></div>
          </div>
        </div>
      </div>
    </div>

    <!-- UI Overlay (Right Panel) -->
    <div
      class="absolute top-0 right-0 p-4 w-80 h-full backdrop-blur-md bg-white/30 transform transition-transform duration-300 flex flex-col">
      <!-- Header with Save Button -->
      <div class="flex items-center justify-between mb-4">
        <h2 class="text-xl font-bold text-gray-800 font-display">Avatar Studio</h2>
        <button @click="saveOutfit" :disabled="isSaving"
          class="px-4 py-2 bg-green-600 hover:bg-green-700 text-white text-sm font-bold rounded-lg shadow-md transform transition-all active:scale-95 disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-2">
          <svg v-if="isSaving" class="animate-spin h-4 w-4" xmlns="http://www.w3.org/2000/svg" fill="none"
            viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4">
            </circle>
            <path class="opacity-75" fill="currentColor"
              d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z">
            </path>
          </svg>
          <span v-else>💾 Save</span>
        </button>
      </div>
      <p v-if="saveMessage" class="text-center mb-2 text-sm font-medium px-2 py-1 rounded-lg"
        :class="saveStatus === 'success' ? 'text-green-700 bg-green-100' : 'text-red-700 bg-red-100'">
        {{ saveMessage }}
      </p>

      <!-- Outfits List (Scrollable, takes remaining space) -->
      <div class="flex-1 overflow-hidden flex flex-col">
        <div class="p-3 bg-white/50 rounded-xl border border-white/40 shadow-sm flex-1 overflow-hidden flex flex-col">
          <h3 class="text-base font-semibold text-gray-700 mb-2 flex items-center gap-2 shrink-0">
            <span class="w-7 h-7 rounded-full bg-green-100 text-green-600 flex items-center justify-center text-sm">
              👕
            </span>
            Recycled Outfits ({{ outfits.length }})
          </h3>
          <div class="flex-1 overflow-y-auto pr-1 custom-scrollbar space-y-1.5">
            <button v-for="outfit in outfits" :key="outfit.id" @click="changeOutfit(outfit.id)"
              class="w-full text-left p-2.5 rounded-lg transition-all duration-200 border-2 flex items-center gap-2.5 group bg-white/80 hover:bg-white"
              :class="currentOutfit === outfit.id ? 'border-green-500 shadow-md transform scale-[1.02]' : 'border-transparent hover:border-gray-200'">
              <div class="w-9 h-9 rounded-full flex items-center justify-center text-base shadow-sm shrink-0"
                :class="outfit.colorClass">
                {{ outfit.icon }}
              </div>
              <div class="min-w-0 flex-1">
                <div class="font-medium text-gray-800 text-sm truncate">{{ outfit.name }}</div>
                <div class="text-xs text-gray-500 truncate">{{ outfit.material }}</div>
              </div>
              <div v-if="currentOutfit === outfit.id" class="text-green-500 shrink-0">
                ✓
              </div>
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount, computed } from 'vue'
import * as THREE from 'three'
import { OrbitControls } from 'three/examples/jsm/controls/OrbitControls.js'
import { useUsersStore } from '@/stores/users'
import gsap from 'gsap'
import { outfits, getOutfitMaterial, getAccessoryConfig } from '@/data/outfits'

// Data
const container = ref<HTMLElement | null>(null)
const currentOutfit = ref('basic')
const usersStore = useUsersStore()
const isSaving = ref(false)
const saveMessage = ref('')
const saveStatus = ref<'success' | 'error' | ''>('')

const currentStats = computed(() => {
  const outfit = outfits.find(o => o.id === currentOutfit.value)
  return outfit ? outfit.stats : []
})

const currentOutfitData = computed(() => {
  return outfits.find(o => o.id === currentOutfit.value)
})

// Three.js variables
let scene: THREE.Scene
let camera: THREE.PerspectiveCamera
let renderer: THREE.WebGLRenderer
let controls: OrbitControls
let avatarGroup: THREE.Group
let environmentGroup: THREE.Group
let head: THREE.Mesh
let body: THREE.Mesh
let leftArm: THREE.Mesh
let rightArm: THREE.Mesh
let leftLeg: THREE.Mesh
let rightLeg: THREE.Mesh
let hatMesh: THREE.Mesh | null = null
let backObjectMesh: THREE.Mesh | null = null
let animationFrameId: number

// Face meshes for animation
let leftEyeWhite: THREE.Mesh
let rightEyeWhite: THREE.Mesh
let leftPupil: THREE.Mesh
let rightPupil: THREE.Mesh
let leftEyebrow: THREE.Mesh
let rightEyebrow: THREE.Mesh
let mouth: THREE.Mesh
let lastBlinkTime = 0
let isBlinking = false

// Setup Scene
const initScene = () => {
  if (!container.value) return

  // Scene
  scene = new THREE.Scene()
  scene.background = new THREE.Color(0xf0fdf4) // Light green background
  scene.fog = new THREE.Fog(0xf0fdf4, 10, 50)

  // Camera
  camera = new THREE.PerspectiveCamera(45, container.value.clientWidth / container.value.clientHeight, 0.1, 1000)
  camera.position.set(0, 1, 5)

  // Renderer
  renderer = new THREE.WebGLRenderer({ antialias: true, alpha: true })
  renderer.setSize(container.value.clientWidth, container.value.clientHeight)
  renderer.shadowMap.enabled = true
  renderer.shadowMap.type = THREE.PCFSoftShadowMap
  renderer.setPixelRatio(Math.min(window.devicePixelRatio, 2))
  container.value.appendChild(renderer.domElement)

  // Controls
  controls = new OrbitControls(camera, renderer.domElement)
  controls.enableDamping = true
  controls.dampingFactor = 0.05
  controls.minDistance = 2
  controls.maxDistance = 10
  controls.maxPolarAngle = Math.PI / 2 - 0.1
  controls.target.set(0, 0, 0)

  // Lighting
  const ambientLight = new THREE.AmbientLight(0xffffff, 0.6)
  scene.add(ambientLight)

  const dirLight = new THREE.DirectionalLight(0xffffff, 1)
  dirLight.position.set(5, 10, 7)
  dirLight.castShadow = true
  dirLight.shadow.mapSize.width = 1024
  dirLight.shadow.mapSize.height = 1024
  scene.add(dirLight)

  const fillLight = new THREE.DirectionalLight(0xa5f3fc, 0.5)
  fillLight.position.set(-5, 0, -5)
  scene.add(fillLight)

  // === ENVIRONMENT ===
  environmentGroup = new THREE.Group()

  // Grass Ground
  const groundGeo = new THREE.PlaneGeometry(20, 20)
  const groundMat = new THREE.MeshStandardMaterial({
    color: 0x4a7c59, // Grass green
    roughness: 0.9,
    metalness: 0.0
  })
  const ground = new THREE.Mesh(groundGeo, groundMat)
  ground.rotation.x = -Math.PI / 2
  ground.position.y = -1.5
  ground.receiveShadow = true
  environmentGroup.add(ground)

  // Walking Path (stone/concrete path)
  const pathGeo = new THREE.PlaneGeometry(4, 20)
  const pathMat = new THREE.MeshStandardMaterial({
    color: 0x9e9e9e, // Gray stone
    roughness: 0.8,
    metalness: 0.1
  })
  const path = new THREE.Mesh(pathGeo, pathMat)
  path.rotation.x = -Math.PI / 2
  path.position.y = -1.49 // Slightly above grass
  path.receiveShadow = true
  environmentGroup.add(path)

  // Path edge lines
  const edgeGeo = new THREE.PlaneGeometry(0.1, 20)
  const edgeMat = new THREE.MeshStandardMaterial({ color: 0x78909c })
  const leftEdge = new THREE.Mesh(edgeGeo, edgeMat)
  leftEdge.rotation.x = -Math.PI / 2
  leftEdge.position.set(-2, -1.48, 0)
  environmentGroup.add(leftEdge)
  const rightEdge = new THREE.Mesh(edgeGeo, edgeMat)
  rightEdge.rotation.x = -Math.PI / 2
  rightEdge.position.set(2, -1.48, 0)
  environmentGroup.add(rightEdge)

  // Helper function to create simple tree
  const createTree = (x: number, z: number, scale = 1) => {
    const treeGroup = new THREE.Group()

    // Trunk
    const trunkGeo = new THREE.CylinderGeometry(0.1 * scale, 0.15 * scale, 0.8 * scale, 8)
    const trunkMat = new THREE.MeshStandardMaterial({ color: 0x5d4037 })
    const trunk = new THREE.Mesh(trunkGeo, trunkMat)
    trunk.position.y = 0.4 * scale
    trunk.castShadow = true
    treeGroup.add(trunk)

    // Foliage (3 spheres stacked)
    const foliageMat = new THREE.MeshStandardMaterial({ color: 0x2e7d32 })
    const foliage1 = new THREE.Mesh(new THREE.SphereGeometry(0.4 * scale, 8, 8), foliageMat)
    foliage1.position.y = 0.9 * scale
    foliage1.castShadow = true
    treeGroup.add(foliage1)

    const foliage2 = new THREE.Mesh(new THREE.SphereGeometry(0.3 * scale, 8, 8), foliageMat)
    foliage2.position.y = 1.3 * scale
    foliage2.castShadow = true
    treeGroup.add(foliage2)

    treeGroup.position.set(x, -1.5, z)
    environmentGroup.add(treeGroup)
  }

  // Add trees
  createTree(-4, -2, 1.2)
  createTree(-3.5, 2, 0.9)
  createTree(4, -1, 1.0)
  createTree(4.5, 3, 1.3)
  createTree(-5, 0, 0.8)

  // Helper function to create flower
  const createFlower = (x: number, z: number, color: number) => {
    const flowerGroup = new THREE.Group()

    // Stem
    const stemGeo = new THREE.CylinderGeometry(0.02, 0.02, 0.2, 6)
    const stemMat = new THREE.MeshStandardMaterial({ color: 0x33691e })
    const stem = new THREE.Mesh(stemGeo, stemMat)
    stem.position.y = 0.1
    flowerGroup.add(stem)

    // Flower head
    const petalMat = new THREE.MeshStandardMaterial({ color })
    const petal = new THREE.Mesh(new THREE.SphereGeometry(0.06, 8, 8), petalMat)
    petal.position.y = 0.22
    flowerGroup.add(petal)

    flowerGroup.position.set(x, -1.5, z)
    environmentGroup.add(flowerGroup)
  }

  // Add flowers along the path
  const flowerColors = [0xe91e63, 0xffeb3b, 0x9c27b0, 0xff5722, 0x03a9f4]
  for (let i = 0; i < 15; i++) {
    const side = Math.random() > 0.5 ? 1 : -1
    const x = (2.3 + Math.random() * 0.5) * side
    const z = (Math.random() - 0.5) * 8
    createFlower(x, z, flowerColors[Math.floor(Math.random() * flowerColors.length)])
  }

  // Recycling Bins
  const createRecycleBin = (x: number, z: number, color: number, _label: string) => {
    const binGroup = new THREE.Group()

    // Bin body
    const binGeo = new THREE.CylinderGeometry(0.15, 0.12, 0.4, 16)
    const binMat = new THREE.MeshStandardMaterial({ color, roughness: 0.3, metalness: 0.5 })
    const bin = new THREE.Mesh(binGeo, binMat)
    bin.position.y = 0.2
    bin.castShadow = true
    binGroup.add(bin)

    // Lid
    const lidGeo = new THREE.CylinderGeometry(0.16, 0.15, 0.05, 16)
    const lidMat = new THREE.MeshStandardMaterial({ color: 0x424242, roughness: 0.4, metalness: 0.6 })
    const lid = new THREE.Mesh(lidGeo, lidMat)
    lid.position.y = 0.42
    lid.castShadow = true
    binGroup.add(lid)

    // Recycle symbol (simple ring)
    const ringGeo = new THREE.TorusGeometry(0.06, 0.015, 8, 16)
    const ringMat = new THREE.MeshStandardMaterial({ color: 0xffffff })
    const ring = new THREE.Mesh(ringGeo, ringMat)
    ring.position.set(0, 0.25, 0.14)
    binGroup.add(ring)

    binGroup.position.set(x, -1.5, z)
    environmentGroup.add(binGroup)
  }

  // Add recycling bins (green, blue, yellow for different waste types)
  createRecycleBin(-3, 4, 0x4caf50, 'organic')   // Green - organic
  createRecycleBin(-2.6, 4, 0x2196f3, 'plastic') // Blue - plastic
  createRecycleBin(-2.2, 4, 0xffeb3b, 'paper')   // Yellow - paper

  // Add environment group to scene
  scene.add(environmentGroup)

  // Avatar
  createAvatar()

  // Animation Loop
  animate()

  // Handle Resize
  window.addEventListener('resize', onWindowResize)
}

const createAvatar = () => {
  avatarGroup = new THREE.Group()

  // Materials
  const skinMaterial = new THREE.MeshStandardMaterial({ color: 0xffd1a9, roughness: 0.5 })
  const outfitMat = new THREE.MeshPhysicalMaterial(outfits[0].materialProp)

  // Head
  const headGeo = new THREE.SphereGeometry(0.35, 32, 32)
  head = new THREE.Mesh(headGeo, skinMaterial)
  head.position.set(0, 1.2, 0)
  head.castShadow = true
  avatarGroup.add(head)

  // === ANIME-STYLE FACE ===

  // Eyes (white sclera)
  const eyeWhiteMat = new THREE.MeshStandardMaterial({ color: 0xffffff })
  const eyeWhiteGeo = new THREE.SphereGeometry(0.08, 16, 16)

  leftEyeWhite = new THREE.Mesh(eyeWhiteGeo, eyeWhiteMat)
  leftEyeWhite.position.set(-0.12, 1.25, 0.28)
  leftEyeWhite.scale.set(1, 1.2, 0.6) // Slightly tall anime-style
  avatarGroup.add(leftEyeWhite)

  rightEyeWhite = new THREE.Mesh(eyeWhiteGeo, eyeWhiteMat)
  rightEyeWhite.position.set(0.12, 1.25, 0.28)
  rightEyeWhite.scale.set(1, 1.2, 0.6)
  avatarGroup.add(rightEyeWhite)

  // Pupils (dark iris + pupil)
  const pupilMat = new THREE.MeshStandardMaterial({ color: 0x2d1b4e }) // Dark purple/black
  const irisGeo = new THREE.SphereGeometry(0.045, 16, 16)

  leftPupil = new THREE.Mesh(irisGeo, pupilMat)
  leftPupil.position.set(-0.12, 1.24, 0.33)
  leftPupil.scale.set(1, 1.1, 0.5)
  avatarGroup.add(leftPupil)

  rightPupil = new THREE.Mesh(irisGeo, pupilMat)
  rightPupil.position.set(0.12, 1.24, 0.33)
  rightPupil.scale.set(1, 1.1, 0.5)
  avatarGroup.add(rightPupil)

  // Eye Highlights (anime sparkle)
  const highlightMat = new THREE.MeshStandardMaterial({ color: 0xffffff, emissive: 0xffffff, emissiveIntensity: 0.3 })
  const highlightGeo = new THREE.SphereGeometry(0.015, 8, 8)

  const leftHighlight = new THREE.Mesh(highlightGeo, highlightMat)
  leftHighlight.position.set(-0.10, 1.27, 0.35)
  avatarGroup.add(leftHighlight)

  const rightHighlight = new THREE.Mesh(highlightGeo, highlightMat)
  rightHighlight.position.set(0.14, 1.27, 0.35)
  avatarGroup.add(rightHighlight)

  // Eyebrows
  const eyebrowMat = new THREE.MeshStandardMaterial({ color: 0x4a3728 }) // Dark brown
  const eyebrowGeo = new THREE.BoxGeometry(0.08, 0.015, 0.02)

  leftEyebrow = new THREE.Mesh(eyebrowGeo, eyebrowMat)
  leftEyebrow.position.set(-0.12, 1.38, 0.30)
  leftEyebrow.rotation.z = 0.15 // Slight angle
  avatarGroup.add(leftEyebrow)

  rightEyebrow = new THREE.Mesh(eyebrowGeo, eyebrowMat)
  rightEyebrow.position.set(0.12, 1.38, 0.30)
  rightEyebrow.rotation.z = -0.15
  avatarGroup.add(rightEyebrow)

  // Mouth (cute smile - half torus)
  const mouthMat = new THREE.MeshStandardMaterial({ color: 0xe88a9a }) // Soft pink
  const mouthGeo = new THREE.TorusGeometry(0.06, 0.015, 8, 16, Math.PI) // Half circle

  mouth = new THREE.Mesh(mouthGeo, mouthMat)
  mouth.position.set(0, 1.08, 0.32)
  mouth.rotation.x = 0.2 // Tilt to face forward
  mouth.rotation.z = Math.PI // Flip to smile shape
  avatarGroup.add(mouth)

  // Blush marks (optional cute detail)
  const blushMat = new THREE.MeshStandardMaterial({ color: 0xffb6c1, transparent: true, opacity: 0.5 })
  const blushGeo = new THREE.CircleGeometry(0.04, 16)

  const leftBlush = new THREE.Mesh(blushGeo, blushMat)
  leftBlush.position.set(-0.22, 1.15, 0.28)
  leftBlush.rotation.y = 0.4
  avatarGroup.add(leftBlush)

  const rightBlush = new THREE.Mesh(blushGeo, blushMat)
  rightBlush.position.set(0.22, 1.15, 0.28)
  rightBlush.rotation.y = -0.4
  avatarGroup.add(rightBlush)

  // Body (Outfit)
  const bodyGeo = new THREE.CylinderGeometry(0.3, 0.25, 1, 32)
  body = new THREE.Mesh(bodyGeo, outfitMat)
  body.position.set(0, 0.5, 0)
  body.castShadow = true
  avatarGroup.add(body)

  // Arms (hanging down by sides)
  const armGeo = new THREE.CapsuleGeometry(0.08, 0.6, 4, 16)

  leftArm = new THREE.Mesh(armGeo, outfitMat)
  leftArm.position.set(-0.35, 0.35, 0) // Lower position, closer to body
  leftArm.rotation.z = 0.15 // Slight angle outward
  leftArm.castShadow = true
  avatarGroup.add(leftArm)

  rightArm = new THREE.Mesh(armGeo, outfitMat)
  rightArm.position.set(0.35, 0.35, 0)
  rightArm.rotation.z = -0.15
  rightArm.castShadow = true
  avatarGroup.add(rightArm)

  // Legs
  const legGeo = new THREE.CapsuleGeometry(0.1, 0.8, 4, 16)
  const legMat = new THREE.MeshPhysicalMaterial({ ...outfits[0].materialProp, color: 0x333333 }) // Pants often dark initially

  leftLeg = new THREE.Mesh(legGeo, legMat)
  leftLeg.position.set(-0.15, -0.4, 0)
  leftLeg.castShadow = true
  avatarGroup.add(leftLeg)

  rightLeg = new THREE.Mesh(legGeo, legMat)
  rightLeg.position.set(0.15, -0.4, 0)
  rightLeg.castShadow = true
  avatarGroup.add(rightLeg)

  // Position avatar group so feet touch the ground (ground is at y = -1.5)
  // Legs bottom is at about y = -0.9 in local space, so offset by -0.6
  avatarGroup.position.y = -0.6

  scene.add(avatarGroup)
}

const changeOutfit = (outfitId: string) => {
  if (currentOutfit.value === outfitId) return // Already selected

  currentOutfit.value = outfitId
  const outfit = outfits.find(o => o.id === outfitId)
  if (!outfit || !avatarGroup) return

  // Animation: Spin + Scale bounce
  const timeline = gsap.timeline()

  // Phase 1: Shrink and spin fast
  timeline.to(avatarGroup.scale, {
    x: 0.1,
    y: 0.1,
    z: 0.1,
    duration: 0.25,
    ease: 'power2.in'
  })
  timeline.to(avatarGroup.rotation, {
    y: avatarGroup.rotation.y + Math.PI * 2,
    duration: 0.25,
    ease: 'power2.in'
  }, '<')

  // Phase 2: Change material and accessories (at the peak of animation)
  timeline.call(() => {
    const newMaterial = new THREE.MeshPhysicalMaterial(getOutfitMaterial(outfit))
    body.material = newMaterial
    leftArm.material = newMaterial
    rightArm.material = newMaterial

    // Remove old accessories
    if (hatMesh) {
      avatarGroup.remove(hatMesh)
      hatMesh = null
    }
    if (backObjectMesh) {
      avatarGroup.remove(backObjectMesh)
      backObjectMesh = null
    }

    // Add new hat
    const hatConfig = getAccessoryConfig(outfit.hat)
    if (hatConfig) {
      const hatGeo = createAccessoryGeometry(hatConfig.type)
      const hatMat = new THREE.MeshStandardMaterial({ color: hatConfig.color })
      hatMesh = new THREE.Mesh(hatGeo, hatMat)
      hatMesh.position.set(...hatConfig.position)
      hatMesh.scale.set(...hatConfig.scale)
      hatMesh.castShadow = true
      avatarGroup.add(hatMesh)
    }

    // Add new backObject
    const backConfig = getAccessoryConfig(outfit.backObject)
    if (backConfig) {
      const backGeo = createAccessoryGeometry(backConfig.type)
      const backMat = new THREE.MeshStandardMaterial({ color: backConfig.color })
      backObjectMesh = new THREE.Mesh(backGeo, backMat)
      backObjectMesh.position.set(...backConfig.position)
      backObjectMesh.scale.set(...backConfig.scale)
      backObjectMesh.castShadow = true
      avatarGroup.add(backObjectMesh)
    }
  })

  // Phase 3: Expand with bounce
  timeline.to(avatarGroup.scale, {
    x: 1.15,
    y: 1.15,
    z: 1.15,
    duration: 0.2,
    ease: 'back.out(2)'
  })

  // Phase 4: Settle to normal size
  timeline.to(avatarGroup.scale, {
    x: 1,
    y: 1,
    z: 1,
    duration: 0.15,
    ease: 'power2.out'
  })
}

// Helper to create geometry based on accessory type
const createAccessoryGeometry = (type: string): THREE.BufferGeometry => {
  switch (type) {
    case 'sphere': return new THREE.SphereGeometry(1, 16, 16)
    case 'box': return new THREE.BoxGeometry(1, 1, 1)
    case 'cylinder': return new THREE.CylinderGeometry(1, 1, 1, 16)
    case 'cone': return new THREE.ConeGeometry(1, 1, 16)
    case 'torus': return new THREE.TorusGeometry(1, 0.3, 8, 16)
    default: return new THREE.SphereGeometry(1, 16, 16)
  }
}

const animate = () => {
  animationFrameId = requestAnimationFrame(animate)

  const time = Date.now() * 0.001

  // Environment slides past avatar (treadmill/infinite scroll effect)
  const slideSpeed = 1.5

  // Slide environment in -Z direction (avatar appears to walk forward)
  if (environmentGroup) {
    // Loop position: slide from +5 to -5, then reset
    environmentGroup.position.z = -(((time * slideSpeed) % 10) - 5)
  }

  // Walking animation intensity
  const walkingIntensity = 1.0

  if (avatarGroup) {
    // Avatar stays in place (x = 0)
    avatarGroup.position.x = 0

    // Bob up and down while walking
    const baseY = -0.6
    avatarGroup.position.y = baseY + Math.abs(Math.sin(time * 4)) * 0.04 * walkingIntensity
  }

  // Head looks forward with slight bob
  if (head) {
    head.rotation.y = Math.sin(time * 0.5) * 0.1
    head.rotation.x = Math.sin(time * 4) * 0.04 * walkingIntensity
    head.rotation.z = Math.sin(time * 2) * 0.02
  }

  // === FACE IDLE ANIMATIONS ===

  // Eye Blinking (random interval 2-5 seconds)
  if (leftEyeWhite && rightEyeWhite) {
    if (!isBlinking && time - lastBlinkTime > 2 + Math.random() * 3) {
      isBlinking = true
      lastBlinkTime = time
    }

    if (isBlinking) {
      const blinkProgress = (time - lastBlinkTime) * 8 // Fast blink
      if (blinkProgress < 1) {
        // Close eyes
        const scaleY = Math.max(0.1, 1.2 - blinkProgress * 1.1)
        leftEyeWhite.scale.y = scaleY
        rightEyeWhite.scale.y = scaleY
        if (leftPupil) leftPupil.scale.y = Math.max(0.1, 1.1 - blinkProgress)
        if (rightPupil) rightPupil.scale.y = Math.max(0.1, 1.1 - blinkProgress)
      } else if (blinkProgress < 2) {
        // Open eyes
        const scaleY = 0.1 + (blinkProgress - 1) * 1.1
        leftEyeWhite.scale.y = Math.min(1.2, scaleY)
        rightEyeWhite.scale.y = Math.min(1.2, scaleY)
        if (leftPupil) leftPupil.scale.y = Math.min(1.1, 0.1 + (blinkProgress - 1))
        if (rightPupil) rightPupil.scale.y = Math.min(1.1, 0.1 + (blinkProgress - 1))
      } else {
        // Reset
        leftEyeWhite.scale.y = 1.2
        rightEyeWhite.scale.y = 1.2
        if (leftPupil) leftPupil.scale.y = 1.1
        if (rightPupil) rightPupil.scale.y = 1.1
        isBlinking = false
      }
    }

    // Subtle eye movement (looking around)
    const eyeLookX = Math.sin(time * 0.3) * 0.02
    const eyeLookY = Math.sin(time * 0.5) * 0.01
    leftEyeWhite.position.x = -0.12 + eyeLookX
    rightEyeWhite.position.x = 0.12 + eyeLookX
    leftEyeWhite.position.y = 1.25 + eyeLookY
    rightEyeWhite.position.y = 1.25 + eyeLookY
    if (leftPupil) {
      leftPupil.position.x = -0.12 + eyeLookX
      leftPupil.position.y = 1.24 + eyeLookY
    }
    if (rightPupil) {
      rightPupil.position.x = 0.12 + eyeLookX
      rightPupil.position.y = 1.24 + eyeLookY
    }
  }

  // Eyebrow subtle movement
  if (leftEyebrow && rightEyebrow) {
    const browMove = Math.sin(time * 0.7) * 0.01
    leftEyebrow.position.y = 1.38 + browMove
    rightEyebrow.position.y = 1.38 + browMove
    // Slight rotation for expression
    leftEyebrow.rotation.z = 0.15 + Math.sin(time * 0.4) * 0.05
    rightEyebrow.rotation.z = -0.15 - Math.sin(time * 0.4) * 0.05
  }

  // Mouth - static smile (no animation)

  // Walking arm swing - natural swing forward/back
  const armSwingAmount = Math.sin(time * 4) * 0.5 * walkingIntensity
  if (leftArm) {
    leftArm.rotation.z = 0.15 // Keep slight outward angle
    leftArm.rotation.x = -armSwingAmount // Forward/back swing
  }
  if (rightArm) {
    rightArm.rotation.z = -0.15
    rightArm.rotation.x = armSwingAmount // Opposite direction
  }

  // Walking leg motion
  const legSwing = Math.sin(time * 4) * 0.45 * walkingIntensity
  if (leftLeg) {
    leftLeg.rotation.x = legSwing
    leftLeg.position.y = -0.4 + Math.abs(Math.sin(time * 4)) * 0.03 * walkingIntensity
  }
  if (rightLeg) {
    rightLeg.rotation.x = -legSwing
    rightLeg.position.y = -0.4 + Math.abs(Math.sin(time * 4 + Math.PI)) * 0.03 * walkingIntensity
  }

  // Body tilt and sway while walking
  if (body) {
    body.rotation.z = Math.sin(time * 4) * 0.04 * walkingIntensity
    body.rotation.x = Math.sin(time * 2) * 0.02 // Slight forward lean
    body.scale.x = 1 + Math.sin(time * 2.5) * 0.01
    body.scale.z = 1 + Math.sin(time * 2.5) * 0.01
  }

  if (controls) controls.update()
  if (renderer && scene && camera) renderer.render(scene, camera)
}

const onWindowResize = () => {
  if (!container.value || !camera || !renderer) return

  camera.aspect = container.value.clientWidth / container.value.clientHeight
  camera.updateProjectionMatrix()
  renderer.setSize(container.value.clientWidth, container.value.clientHeight)
}

const saveOutfit = async () => {
  isSaving.value = true
  saveMessage.value = ''
  saveStatus.value = ''

  try {
    const success = await usersStore.updateAvatar(currentOutfit.value)
    if (success) {
      saveMessage.value = 'Outfit saved successfully!'
      saveStatus.value = 'success'
      setTimeout(() => {
        saveMessage.value = ''
        saveStatus.value = ''
      }, 3000)
    } else {
      throw new Error('Failed to save')
    }
  } catch (error) {
    console.error('Save failed:', error)
    saveMessage.value = 'Failed to save outfit.'
    saveStatus.value = 'error'
  } finally {
    isSaving.value = false
  }
}

onMounted(async () => {
  await usersStore.checkLogin()
  const savedOutfit = usersStore.user?.avatar_config
  if (savedOutfit) {
    currentOutfit.value = savedOutfit
  }
  initScene()
  // Apply saved outfit to 3D model after scene is initialized
  if (savedOutfit && body) {
    applyOutfitMaterial(savedOutfit)
  }
})

const applyOutfitMaterial = (outfitId: string) => {
  const outfit = outfits.find(o => o.id === outfitId)
  if (!outfit || !body) return

  const newMaterial = new THREE.MeshPhysicalMaterial(getOutfitMaterial(outfit))
  body.material = newMaterial
  leftArm.material = newMaterial
  rightArm.material = newMaterial

  // Add hat if outfit has one
  const hatConfig = getAccessoryConfig(outfit.hat)
  if (hatConfig && avatarGroup) {
    const hatGeo = createAccessoryGeometry(hatConfig.type)
    const hatMat = new THREE.MeshStandardMaterial({ color: hatConfig.color })
    hatMesh = new THREE.Mesh(hatGeo, hatMat)
    hatMesh.position.set(...hatConfig.position)
    hatMesh.scale.set(...hatConfig.scale)
    hatMesh.castShadow = true
    avatarGroup.add(hatMesh)
  }

  // Add backObject if outfit has one
  const backConfig = getAccessoryConfig(outfit.backObject)
  if (backConfig && avatarGroup) {
    const backGeo = createAccessoryGeometry(backConfig.type)
    const backMat = new THREE.MeshStandardMaterial({ color: backConfig.color })
    backObjectMesh = new THREE.Mesh(backGeo, backMat)
    backObjectMesh.position.set(...backConfig.position)
    backObjectMesh.scale.set(...backConfig.scale)
    backObjectMesh.castShadow = true
    avatarGroup.add(backObjectMesh)
  }
}

onBeforeUnmount(() => {
  window.removeEventListener('resize', onWindowResize)
  cancelAnimationFrame(animationFrameId)
  if (renderer) {
    renderer.dispose()
    const canvas = renderer.domElement
    if (canvas && canvas.parentNode) canvas.parentNode.removeChild(canvas)
  }
})
</script>

<style scoped>
.font-display {
  font-family: 'Outfit', sans-serif;
}

.custom-scrollbar::-webkit-scrollbar {
  width: 6px;
}

.custom-scrollbar::-webkit-scrollbar-track {
  background: transparent;
}

.custom-scrollbar::-webkit-scrollbar-thumb {
  background-color: rgba(156, 163, 175, 0.5);
  border-radius: 20px;
}

.custom-scrollbar::-webkit-scrollbar-thumb:hover {
  background-color: rgba(107, 114, 128, 0.8);
}
</style>
