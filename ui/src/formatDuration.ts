 

export const formatDuration = (ms: number) => {
  const time = {
    year: Math.floor(ms / 1000 / 60 / 60 / 24 / 365),
    day: Math.floor(ms / 1000 / 60 / 60 / 24) % 365,
    hour: Math.floor(ms / 1000 / 60 / 60) % 24,
    minute: Math.floor(ms / 1000 / 60) % 60,
    second: Math.floor(ms / 1000) % 60,
    millisecond: Math.floor(ms) % 1000,
  }
  return Object.entries(time)
    .filter(val => val[1] !== 0)
    .map(([key, val]) => `${val} ${key}${val !== 1 ? 's' : ''}`)
    .join(' ')
}

