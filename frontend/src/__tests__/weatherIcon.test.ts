import { getWeatherIcon } from "@/lib/weatherIcon";

describe("getWeatherIcon", () => {
  it.each([
    ["Snow", "❄️"],
    ["Heavy Snow", "❄️"],
    ["Blizzard", "❄️"],
    ["Sleet", "❄️"],
    ["Ice Pellets", "❄️"],
  ])("returns snowflake for '%s'", (input, expected) => {
    expect(getWeatherIcon(input)).toBe(expected);
  });

  it.each([
    ["Rain", "🌧️"],
    ["Light Rain", "🌧️"],
    ["Showers", "🌧️"],
    ["Chance Rain Showers", "🌧️"],
    ["Drizzle", "🌧️"],
    ["Thunderstorm", "🌧️"],
    ["Scattered Thunderstorms", "🌧️"],
  ])("returns rain for '%s'", (input, expected) => {
    expect(getWeatherIcon(input)).toBe(expected);
  });

  it.each([
    ["Cloudy", "☁️"],
    ["Mostly Cloudy", "☁️"],
    ["Overcast", "☁️"],
    ["Fog", "☁️"],
    ["Dense Fog", "☁️"],
  ])("returns cloud for '%s'", (input, expected) => {
    expect(getWeatherIcon(input)).toBe(expected);
  });

  it.each([
    ["Partly Sunny", "⛅"],
    ["Partly Clear", "⛅"],
  ])("returns partly cloudy for '%s'", (input, expected) => {
    expect(getWeatherIcon(input)).toBe(expected);
  });

  it.each([
    ["Sunny", "☀️"],
    ["Clear", "☀️"],
    ["Hot", "☀️"],
    ["", "☀️"],
  ])("returns sun for '%s'", (input, expected) => {
    expect(getWeatherIcon(input)).toBe(expected);
  });

  it("is case-insensitive", () => {
    expect(getWeatherIcon("HEAVY SNOW")).toBe("❄️");
    expect(getWeatherIcon("light rain")).toBe("🌧️");
    expect(getWeatherIcon("CLOUDY")).toBe("☁️");
  });

  it("prioritizes snow over rain for mixed conditions", () => {
    expect(getWeatherIcon("Rain and Snow")).toBe("❄️");
  });

  it("prioritizes rain over cloud for mixed conditions", () => {
    expect(getWeatherIcon("Rain and Cloudy")).toBe("🌧️");
  });
});
