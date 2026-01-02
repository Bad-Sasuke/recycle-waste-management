import * as tf from '@tensorflow/tfjs'
import * as mobilenet from '@tensorflow-models/mobilenet'

// Define waste categories and keywords for simple mapping
const CATEGORY_KEYWORDS: Record<string, string[]> = {
  Plastic: [
    'bottle',
    'tub',
    'crate',
    'nipple',
    'bucket',
    'plastic',
    'container',
    'cup',
    'water bottle',
    'soda bottle',
  ],
  Glass: ['bottle', 'goblet', 'flask', 'beaker', 'jar', 'glass', 'wine bottle', 'beer bottle'],
  Paper: [
    'carton',
    'envelope',
    'menu',
    'packet',
    'tissue',
    'box',
    'cardboard',
    'paper',
    'book',
    'newspaper',
  ],
  Metal: ['can', 'opener', 'nail', 'screw', 'metal', 'aluminum', 'tin', 'steel'],
  Electronics: [
    'monitor',
    'screen',
    'keyboard',
    'mouse',
    'laptop',
    'phone',
    'computer',
    'electronic',
  ],
}

export interface WastePrediction {
  className: string
  probability: number
  predictedCategory?: string
}

class WasteClassifierService {
  private model: mobilenet.MobileNet | null = null
  private isLoading = false

  async loadModel() {
    if (this.model || this.isLoading) return
    this.isLoading = true
    try {
      console.log('Loading MobileNet model...')
      await tf.ready()
      this.model = await mobilenet.load({
        version: 2,
        alpha: 1.0,
      })
      console.log('MobileNet model loaded')
    } catch (error) {
      console.error('Failed to load model:', error)
      throw error
    } finally {
      this.isLoading = false
    }
  }

  async classify(
    img: HTMLImageElement | HTMLVideoElement | HTMLCanvasElement,
  ): Promise<WastePrediction[]> {
    if (!this.model) {
      await this.loadModel()
    }

    if (!this.model) {
      throw new Error('Model not loaded')
    }

    const predictions = await this.model.classify(img)

    return predictions.map((p) => ({
      className: p.className,
      probability: p.probability,
      predictedCategory: this.mapToCategory(p.className),
    }))
  }

  private mapToCategory(className: string): string | undefined {
    const lowerName = className.toLowerCase()
    for (const [category, keywords] of Object.entries(CATEGORY_KEYWORDS)) {
      if (keywords.some((k) => lowerName.includes(k))) {
        // Special case for bottles which can be glass or plastic
        if (lowerName.includes('bottle')) {
          if (
            lowerName.includes('water') ||
            lowerName.includes('soda') ||
            lowerName.includes('pop')
          ) {
            return 'Plastic'
          }
          if (lowerName.includes('beer') || lowerName.includes('wine')) {
            return 'Glass'
          }
          // Default to Plastic if unsure
          return 'Plastic'
        }
        return category
      }
    }
    return undefined
  }
}

export const wasteClassifier = new WasteClassifierService()
