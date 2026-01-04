export interface CandleData {
  x: number // timestamp
  y: [number, number, number, number] // [Open, High, Low, Close]
}

export interface OrderBookItem {
  price: number
  volume: number
  total: number
  percent: number // for visual bar background
}

export interface MarketIndex {
  name: string
  value: number
  change: number
  history: number[]
}

// Generate realistic-looking OHLC data
// y: [Open, High, Low, Close]
export function generateCandles(
  startPrice: number,
  count: number,
  intervalMs: number = 60000,
): CandleData[] {
  const candles: CandleData[] = []
  let currentPrice = startPrice
  let currentTime = Date.now() - count * intervalMs

  for (let i = 0; i < count; i++) {
    const volatility = currentPrice * 0.05 // 5% volatility for dramatic effect
    const change = (Math.random() - 0.5) * volatility

    let open = currentPrice
    let close = currentPrice + change

    // Ensure reasonable positive prices
    if (open < 1) open = 1
    if (close < 1) close = 1

    // Low must be lower than min(open, close), High must be higher than max(open, close)
    const minBody = Math.min(open, close)
    const maxBody = Math.max(open, close)

    // Random wicks
    const low = minBody - Math.random() * volatility * 0.2
    const high = maxBody + Math.random() * volatility * 0.2

    candles.push({
      x: currentTime,
      y: [
        parseFloat(open.toFixed(2)),
        parseFloat(high.toFixed(2)),
        parseFloat(low.toFixed(2)),
        parseFloat(close.toFixed(2)),
      ],
    })

    currentPrice = close
    currentTime += intervalMs
  }

  return candles
}

// Generate Order Book (Bids/Asks)
export function generateOrderBook(currentPrice: number, depth: number = 10) {
  const bids: OrderBookItem[] = []
  const asks: OrderBookItem[] = []

  let bidPrice = currentPrice
  let askPrice = currentPrice

  for (let i = 0; i < depth; i++) {
    // Bids go down
    bidPrice -= Math.random() * currentPrice * 0.005
    if (bidPrice < 0) bidPrice = 0.1
    const bidVol = Math.floor(Math.random() * 500) + 10
    bids.push({
      price: parseFloat(bidPrice.toFixed(2)),
      volume: bidVol,
      total: 0,
      percent: 0,
    })

    // Asks go up
    askPrice += Math.random() * currentPrice * 0.005
    const askVol = Math.floor(Math.random() * 500) + 10
    asks.push({
      price: parseFloat(askPrice.toFixed(2)),
      volume: askVol,
      total: 0,
      percent: 0,
    })
  }

  // Calculate percentages for visual bars
  const maxVol = Math.max(...bids.map((b) => b.volume), ...asks.map((a) => a.volume))
  bids.forEach((b) => (b.percent = (b.volume / maxVol) * 100))
  asks.forEach((a) => (a.percent = (a.volume / maxVol) * 100))

  return { bids, asks }
}

// Generate Correlation Matrix
export function generateCorrelationMatrix(assets: string[]) {
  const matrix: { x: string; y: string; value: number }[] = []

  for (let i = 0; i < assets.length; i++) {
    for (let j = 0; j < assets.length; j++) {
      let value = 0
      if (i === j) {
        value = 1
      } else {
        // Random correlation between -1 and 1
        value = Math.random() * 2 - 1
      }
      matrix.push({
        x: assets[i],
        y: assets[j],
        value: parseFloat(value.toFixed(2)),
      })
    }
  }
  return matrix
}

// Generate Volume Oscillator Data
export function generateVolumeData(count: number, intervalMs: number = 60000) {
  const data: { x: number; y: number }[] = []
  let currentTime = Date.now() - count * intervalMs

  for (let i = 0; i < count; i++) {
    // Random high volume
    const volume = Math.floor(Math.random() * 1000) + 100
    data.push({
      x: currentTime,
      y: volume,
    })
    currentTime += intervalMs
  }
  return data
}
