import * as THREE from 'three'
import outfitsData from './outfits.json'

// Types
export interface AccessoryConfig {
  type: 'sphere' | 'box' | 'cylinder' | 'cone' | 'torus' | 'none'
  color: string
  scale?: [number, number, number]
  position?: [number, number, number]
}

export interface OutfitData {
  id: string
  name: string
  material: string
  icon: string
  colorClass: string
  stats: { label: string; value: number }[]
  materialProp: Record<string, unknown>
  hat?: AccessoryConfig
  backObject?: AccessoryConfig
}

// Parse color string to number
const parseColor = (colorStr: string): number => {
  if (colorStr.startsWith('0x')) {
    return parseInt(colorStr, 16)
  }
  return parseInt(colorStr.replace('#', ''), 16)
}

// Convert JSON outfit to Three.js compatible format
export const getOutfitMaterial = (outfit: OutfitData): THREE.MeshPhysicalMaterialParameters => {
  const props: THREE.MeshPhysicalMaterialParameters = {}

  for (const [key, value] of Object.entries(outfit.materialProp)) {
    if (key === 'color' || key === 'emissive') {
      ;(props as Record<string, unknown>)[key] = parseColor(value as string)
    } else {
      ;(props as Record<string, unknown>)[key] = value
    }
  }

  return props
}

// Get accessory config with parsed colors
export const getAccessoryConfig = (config?: AccessoryConfig) => {
  if (!config || config.type === 'none') return null

  return {
    type: config.type,
    color: parseColor(config.color),
    scale: config.scale || [0.2, 0.2, 0.2],
    position: config.position || [0, 1.5, 0],
  }
}

// Export raw data
export const outfits: OutfitData[] = outfitsData as OutfitData[]

// Get outfit by ID
export const getOutfitById = (id: string): OutfitData | undefined => {
  return outfits.find((o) => o.id === id)
}
