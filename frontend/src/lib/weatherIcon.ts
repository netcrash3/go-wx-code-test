enum WeatherKeyword {
  Snow = "snow",
  Blizzard = "blizzard",
  Sleet = "sleet",
  Ice = "ice",
  Rain = "rain",
  Shower = "shower",
  Drizzle = "drizzle",
  Thunderstorm = "thunderstorm",
  Cloud = "cloud",
  Overcast = "overcast",
  Fog = "fog",
  Partly = "partly",
}

const SNOW_ICON = "❄️";
const RAIN_ICON = "🌧️";
const CLOUD_ICON = "☁️";
const PARTLY_CLOUDY_ICON = "⛅";
const SUNNY_ICON = "☀️";

const SNOW_KEYWORDS: WeatherKeyword[] = [
  WeatherKeyword.Snow,
  WeatherKeyword.Blizzard,
  WeatherKeyword.Sleet,
  WeatherKeyword.Ice,
];

const RAIN_KEYWORDS: WeatherKeyword[] = [
  WeatherKeyword.Rain,
  WeatherKeyword.Shower,
  WeatherKeyword.Drizzle,
  WeatherKeyword.Thunderstorm,
];

const CLOUD_KEYWORDS: WeatherKeyword[] = [
  WeatherKeyword.Cloud,
  WeatherKeyword.Overcast,
  WeatherKeyword.Fog,
];

export function getWeatherIcon(shortForecast: string): string {
  const lower = shortForecast.toLowerCase();

  if (SNOW_KEYWORDS.some((keyword) => lower.includes(keyword))) {
    return SNOW_ICON;
  }
  if (RAIN_KEYWORDS.some((keyword) => lower.includes(keyword))) {
    return RAIN_ICON;
  }
  if (CLOUD_KEYWORDS.some((keyword) => lower.includes(keyword))) {
    return CLOUD_ICON;
  }
  if (lower.includes(WeatherKeyword.Partly)) {
    return PARTLY_CLOUDY_ICON;
  }
  return SUNNY_ICON;
}
