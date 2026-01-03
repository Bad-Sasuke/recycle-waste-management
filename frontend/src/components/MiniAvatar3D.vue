<template>
    <div ref="container" class="w-full h-full relative cursor-grab active:cursor-grabbing"></div>
</template>

<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount, watch } from 'vue'
import * as THREE from 'three'
import { OrbitControls } from 'three/examples/jsm/controls/OrbitControls.js'
import { outfits, getOutfitMaterial, getAccessoryConfig, type AccessoryConfig } from '@/data/outfits'

const props = defineProps<{
    outfitId?: string
}>()

const container = ref<HTMLElement | null>(null)
let scene: THREE.Scene
let camera: THREE.PerspectiveCamera
let renderer: THREE.WebGLRenderer
let controls: OrbitControls
let avatarGroup: THREE.Group
let animationFrameId: number

// Mesh references for updates
let body: THREE.Mesh
let leftArm: THREE.Mesh
let rightArm: THREE.Mesh
let hatMesh: THREE.Mesh | null = null
let backObjectMesh: THREE.Mesh | null = null

// Face animation refs
let leftEyeWhite: THREE.Mesh, rightEyeWhite: THREE.Mesh
let leftPupil: THREE.Mesh, rightPupil: THREE.Mesh
let leftEyebrow: THREE.Mesh, rightEyebrow: THREE.Mesh
let lastBlinkTime = 0
let isBlinking = false

const initScene = () => {
    if (!container.value) return

    // 1. Scene setup (Transparent BG)
    scene = new THREE.Scene()
    // No background color set = transparent

    // 2. Camera setup for closer look
    const aspect = container.value.clientWidth / container.value.clientHeight
    camera = new THREE.PerspectiveCamera(45, aspect, 0.1, 100)
    camera.position.set(0, 1, 3.5) // Closer than original

    // 3. Renderer
    renderer = new THREE.WebGLRenderer({ antialias: true, alpha: true })
    renderer.setSize(container.value.clientWidth, container.value.clientHeight)
    renderer.shadowMap.enabled = true
    renderer.shadowMap.type = THREE.PCFSoftShadowMap
    renderer.setPixelRatio(Math.min(window.devicePixelRatio, 2))
    container.value.appendChild(renderer.domElement)

    // 4. Controls
    controls = new OrbitControls(camera, renderer.domElement)
    controls.enableDamping = true
    controls.dampingFactor = 0.05
    controls.enableZoom = false
    controls.enablePan = false
    controls.minPolarAngle = Math.PI / 3
    controls.maxPolarAngle = Math.PI / 2
    controls.autoRotate = false

    // 5. Lights
    const ambientLight = new THREE.AmbientLight(0xffffff, 0.7)
    scene.add(ambientLight)

    const dirLight = new THREE.DirectionalLight(0xffffff, 1.2)
    dirLight.position.set(2, 5, 3)
    dirLight.castShadow = true
    scene.add(dirLight)

    const fillLight = new THREE.DirectionalLight(0xa5f3fc, 0.5)
    fillLight.position.set(-2, 0, -2)
    scene.add(fillLight)

    // 6. Create Avatar
    createAvatar()

    // 7. Apply props if exists
    if (props.outfitId) {
        updateOutfit(props.outfitId)
    }

    // 8. Animation Loop
    animate()

    // 9. Resize handler
    window.addEventListener('resize', onWindowResize)
}

const createAvatar = () => {
    avatarGroup = new THREE.Group()

    // Materials
    const skinMaterial = new THREE.MeshStandardMaterial({ color: 0xffd1a9, roughness: 0.5 })
    // Default outfit
    const defaultOutfit = outfits.find(o => o.id === (props.outfitId || 'basic')) || outfits[0]
    const outfitMat = new THREE.MeshPhysicalMaterial(defaultOutfit.materialProp)

    // Head
    const head = new THREE.Mesh(new THREE.SphereGeometry(0.35, 32, 32), skinMaterial)
    head.position.set(0, 1.2, 0)
    head.castShadow = true
    avatarGroup.add(head)

    // === Anime Face Details (Simplified creation) ===
    const eyeWhiteGeo = new THREE.SphereGeometry(0.08, 16, 16)
    const eyeWhiteMat = new THREE.MeshStandardMaterial({ color: 0xffffff })
    leftEyeWhite = new THREE.Mesh(eyeWhiteGeo, eyeWhiteMat)
    rightEyeWhite = new THREE.Mesh(eyeWhiteGeo, eyeWhiteMat)

    leftEyeWhite.position.set(-0.12, 1.25, 0.28); leftEyeWhite.scale.set(1, 1.2, 0.6)
    rightEyeWhite.position.set(0.12, 1.25, 0.28); rightEyeWhite.scale.set(1, 1.2, 0.6)
    avatarGroup.add(leftEyeWhite); avatarGroup.add(rightEyeWhite)

    const pupilMat = new THREE.MeshStandardMaterial({ color: 0x2d1b4e })
    const irisGeo = new THREE.SphereGeometry(0.045, 16, 16)
    leftPupil = new THREE.Mesh(irisGeo, pupilMat)
    rightPupil = new THREE.Mesh(irisGeo, pupilMat)
    leftPupil.position.set(-0.12, 1.24, 0.33); leftPupil.scale.set(1, 1.1, 0.5)
    rightPupil.position.set(0.12, 1.24, 0.33); rightPupil.scale.set(1, 1.1, 0.5)
    avatarGroup.add(leftPupil); avatarGroup.add(rightPupil)

    // Eyebrows
    const eyebrowGeo = new THREE.BoxGeometry(0.08, 0.015, 0.02)
    const eyebrowMat = new THREE.MeshStandardMaterial({ color: 0x4a3728 })
    leftEyebrow = new THREE.Mesh(eyebrowGeo, eyebrowMat)
    rightEyebrow = new THREE.Mesh(eyebrowGeo, eyebrowMat)
    leftEyebrow.position.set(-0.12, 1.38, 0.30); leftEyebrow.rotation.z = 0.15
    rightEyebrow.position.set(0.12, 1.38, 0.30); rightEyebrow.rotation.z = -0.15
    avatarGroup.add(leftEyebrow); avatarGroup.add(rightEyebrow)

    // Mouth
    const mouthMat = new THREE.MeshStandardMaterial({ color: 0xe88a9a })
    const mouth = new THREE.Mesh(new THREE.TorusGeometry(0.06, 0.015, 8, 16, Math.PI), mouthMat)
    mouth.position.set(0, 1.08, 0.32); mouth.rotation.z = Math.PI; mouth.rotation.x = 0.2
    avatarGroup.add(mouth)

    // Body
    body = new THREE.Mesh(new THREE.CylinderGeometry(0.3, 0.25, 1, 32), outfitMat)
    body.position.set(0, 0.5, 0); body.castShadow = true
    avatarGroup.add(body)

    // Arms
    const armGeo = new THREE.CapsuleGeometry(0.08, 0.6, 4, 16)
    leftArm = new THREE.Mesh(armGeo, outfitMat)
    leftArm.position.set(-0.35, 0.35, 0); leftArm.rotation.z = 0.15; leftArm.castShadow = true
    rightArm = new THREE.Mesh(armGeo, outfitMat)
    rightArm.position.set(0.35, 0.35, 0); rightArm.rotation.z = -0.15; rightArm.castShadow = true
    avatarGroup.add(leftArm); avatarGroup.add(rightArm)

    // Legs (Dark pants default)
    const legGeo = new THREE.CapsuleGeometry(0.1, 0.8, 4, 16)
    const legMat = new THREE.MeshPhysicalMaterial({ ...defaultOutfit.materialProp, color: 0x333333 })
    const leftLeg = new THREE.Mesh(legGeo, legMat)
    leftLeg.name = 'leftLeg' // For animation
    leftLeg.position.set(-0.15, -0.4, 0); leftLeg.castShadow = true
    const rightLeg = new THREE.Mesh(legGeo, legMat)
    rightLeg.name = 'rightLeg'
    rightLeg.position.set(0.15, -0.4, 0); rightLeg.castShadow = true
    avatarGroup.add(leftLeg); avatarGroup.add(rightLeg)

    // Center Y
    avatarGroup.position.y = -0.6

    // Initial accessories
    if (defaultOutfit.hat) addAccessory(defaultOutfit.hat, 'hat')
    if (defaultOutfit.backObject) addAccessory(defaultOutfit.backObject, 'back')

    scene.add(avatarGroup)
}

const addAccessory = (configData: AccessoryConfig, slot: 'hat' | 'back') => {
    const config = getAccessoryConfig(configData)
    if (!config) return

    const geo = createAccessoryGeometry(config.type)
    const mat = new THREE.MeshStandardMaterial({ color: config.color })
    const mesh = new THREE.Mesh(geo, mat)
    mesh.position.set(...config.position)
    mesh.scale.set(...config.scale)
    mesh.castShadow = true

    avatarGroup.add(mesh)

    if (slot === 'hat') hatMesh = mesh
    else backObjectMesh = mesh
}

const createAccessoryGeometry = (type: string) => {
    switch (type) {
        case 'sphere': return new THREE.SphereGeometry(1, 16, 16)
        case 'box': return new THREE.BoxGeometry(1, 1, 1)
        case 'cylinder': return new THREE.CylinderGeometry(1, 1, 1, 16)
        case 'cone': return new THREE.ConeGeometry(1, 1, 16)
        case 'torus': return new THREE.TorusGeometry(1, 0.3, 8, 16)
        default: return new THREE.SphereGeometry(1, 16, 16)
    }
}

const updateOutfit = (newOutfitId: string) => {
    const outfit = outfits.find(o => o.id === newOutfitId)
    if (!outfit || !body) return

    const newMat = new THREE.MeshPhysicalMaterial(getOutfitMaterial(outfit))
    body.material = newMat
    leftArm.material = newMat
    rightArm.material = newMat

    // Remove old acc
    if (hatMesh) { avatarGroup.remove(hatMesh); hatMesh = null }
    if (backObjectMesh) { avatarGroup.remove(backObjectMesh); backObjectMesh = null }

    // Add new acc
    if (outfit.hat) addAccessory(outfit.hat, 'hat')
    if (outfit.backObject) addAccessory(outfit.backObject, 'back')
}

watch(() => props.outfitId, (newId) => {
    if (newId) updateOutfit(newId)
})

const animate = () => {
    animationFrameId = requestAnimationFrame(animate)
    const time = Date.now() * 0.001

    // Walking Animation
    if (avatarGroup) {
        // Bob
        avatarGroup.position.y = -0.6 + Math.abs(Math.sin(time * 4)) * 0.04
        // Walk Rotation
        avatarGroup.rotation.y = Math.sin(time * 0.5) * 0.2 // Look around slowly
    }

    // Arms
    if (leftArm) leftArm.rotation.x = -Math.sin(time * 4) * 0.5
    if (rightArm) rightArm.rotation.x = Math.sin(time * 4) * 0.5

    // Legs
    const leftLeg = avatarGroup?.getObjectByName('leftLeg')
    const rightLeg = avatarGroup?.getObjectByName('rightLeg')
    if (leftLeg) {
        leftLeg.rotation.x = Math.sin(time * 4) * 0.45
        leftLeg.position.y = -0.4 + Math.abs(Math.sin(time * 4)) * 0.03
    }
    if (rightLeg) {
        rightLeg.rotation.x = -Math.sin(time * 4) * 0.45
        rightLeg.position.y = -0.4 + Math.abs(Math.sin(time * 4 + Math.PI)) * 0.03
    }

    // Blinking logic same as before...
    if (!isBlinking && time - lastBlinkTime > 2 + Math.random() * 3) {
        isBlinking = true
        lastBlinkTime = time
    }
    if (isBlinking) {
        const blinkProgress = (time - lastBlinkTime) * 8
        if (blinkProgress < 1) {
            const scaleY = Math.max(0.1, 1.2 - blinkProgress * 1.1)
            leftEyeWhite.scale.y = scaleY; rightEyeWhite.scale.y = scaleY
        } else if (blinkProgress < 2) {
            const scaleY = 0.1 + (blinkProgress - 1) * 1.1
            leftEyeWhite.scale.y = Math.min(1.2, scaleY); rightEyeWhite.scale.y = Math.min(1.2, scaleY)
        } else {
            leftEyeWhite.scale.y = 1.2; rightEyeWhite.scale.y = 1.2; isBlinking = false
        }
    }

    controls.update()
    renderer.render(scene, camera)
}

const onWindowResize = () => {
    if (!container.value) return
    camera.aspect = container.value.clientWidth / container.value.clientHeight
    camera.updateProjectionMatrix()
    renderer.setSize(container.value.clientWidth, container.value.clientHeight)
}

onMounted(() => {
    initScene()
})

onBeforeUnmount(() => {
    window.removeEventListener('resize', onWindowResize)
    cancelAnimationFrame(animationFrameId)
    if (renderer) renderer.dispose()
})
</script>
