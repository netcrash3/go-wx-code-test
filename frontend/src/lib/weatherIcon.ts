export function getWeatherIcon(shortForecast: string): string {
  const lower = shortForecast.toLowerCase();
  if (lower.includes("snow") || lower.includes("blizzard") || lower.includes("sleet") || lower.includes("ice")) {
    return "❄️";
  }
  if (lower.includes("rain") || lower.includes("shower") || lower.includes("drizzle") || lower.includes("thunderstorm")) {
    return "🌧️";
  }
  if (lower.includes("cloud") || lower.includes("overcast") || lower.includes("fog")) {
    return "☁️";
  }
  if (lower.includes("partly")) {
    return "⛅";
  }
  return "☀️";
}
