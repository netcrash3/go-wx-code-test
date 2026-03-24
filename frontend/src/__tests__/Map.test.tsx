import { render, screen, act, waitFor } from "@testing-library/react";

// Mock react-leaflet before importing the component
const mockFlyTo = jest.fn();

jest.mock("react-leaflet", () => ({
  MapContainer: ({ children }: { children: React.ReactNode }) => (
    <div data-testid="map-container">{children}</div>
  ),
  TileLayer: () => <div data-testid="tile-layer" />,
  Marker: ({
    children,
    position,
  }: {
    children: React.ReactNode;
    position: [number, number];
  }) => (
    <div data-testid="marker" data-position={position.join(",")}>
      {children}
    </div>
  ),
  Tooltip: ({ children }: { children: React.ReactNode }) => (
    <div data-testid="tooltip">{children}</div>
  ),
  useMap: () => ({ flyTo: mockFlyTo }),
  useMapEvents: () => null,
}));

jest.mock("leaflet", () => ({
  Icon: { Default: { mergeOptions: jest.fn() } },
}));

jest.mock("leaflet/dist/leaflet.css", () => ({}));
jest.mock("leaflet/dist/images/marker-icon-2x.png", () => ({ src: "" }));
jest.mock("leaflet/dist/images/marker-icon.png", () => ({ src: "" }));
jest.mock("leaflet/dist/images/marker-shadow.png", () => ({ src: "" }));

import Map from "@/components/Map";

describe("Map component", () => {
  let originalGeolocation: Geolocation;

  beforeEach(() => {
    originalGeolocation = navigator.geolocation;
    jest.resetAllMocks();
    global.fetch = jest.fn();
  });

  afterEach(() => {
    Object.defineProperty(navigator, "geolocation", {
      value: originalGeolocation,
      writable: true,
    });
  });

  it("renders the map container and tile layer", () => {
    Object.defineProperty(navigator, "geolocation", {
      value: undefined,
      writable: true,
    });

    render(<Map />);

    expect(screen.getByTestId("map-container")).toBeInTheDocument();
    expect(screen.getByTestId("tile-layer")).toBeInTheDocument();
  });

  it("does not render a marker when geolocation is unavailable", () => {
    Object.defineProperty(navigator, "geolocation", {
      value: undefined,
      writable: true,
    });

    render(<Map />);

    expect(screen.queryByTestId("marker")).not.toBeInTheDocument();
  });

  it("renders a marker when geolocation succeeds", async () => {
    const mockPosition = {
      coords: { latitude: 40.7128, longitude: -74.006 },
    };

    Object.defineProperty(navigator, "geolocation", {
      value: {
        getCurrentPosition: (success: PositionCallback) =>
          success(mockPosition as GeolocationPosition),
      },
      writable: true,
    });

    (global.fetch as jest.Mock).mockResolvedValue({
      ok: true,
      json: () =>
        Promise.resolve({
          shortForecast: "Sunny",
          temperature: 72,
          temperatureUnit: "F",
          temperatureCharacterization: "moderate",
        }),
    });

    await act(async () => {
      render(<Map />);
    });

    expect(screen.getByTestId("marker")).toBeInTheDocument();
    expect(screen.getByTestId("marker")).toHaveAttribute(
      "data-position",
      "40.7128,-74.006"
    );
  });

  it("shows spinner in tooltip while forecast is loading", async () => {
    let resolvePromise!: (value: unknown) => void;
    const pendingPromise = new Promise((resolve) => {
      resolvePromise = resolve;
    });

    Object.defineProperty(navigator, "geolocation", {
      value: {
        getCurrentPosition: (success: PositionCallback) =>
          success({
            coords: { latitude: 40, longitude: -74 },
          } as GeolocationPosition),
      },
      writable: true,
    });

    (global.fetch as jest.Mock).mockReturnValue(pendingPromise);

    await act(async () => {
      render(<Map />);
    });

    const tooltip = screen.getByTestId("tooltip");
    expect(tooltip.querySelector("svg")).toBeInTheDocument();

    // Clean up pending promise
    resolvePromise({
      ok: true,
      json: () =>
        Promise.resolve({
          shortForecast: "Clear",
          temperature: 65,
          temperatureUnit: "F",
          temperatureCharacterization: "moderate",
        }),
    });
  });

  it("displays forecast data once loaded", async () => {
    Object.defineProperty(navigator, "geolocation", {
      value: {
        getCurrentPosition: (success: PositionCallback) =>
          success({
            coords: { latitude: 40, longitude: -74 },
          } as GeolocationPosition),
      },
      writable: true,
    });

    (global.fetch as jest.Mock).mockResolvedValue({
      ok: true,
      json: () =>
        Promise.resolve({
          shortForecast: "Partly Cloudy",
          temperature: 68,
          temperatureUnit: "F",
          temperatureCharacterization: "moderate",
        }),
    });

    await act(async () => {
      render(<Map />);
    });

    await waitFor(() => {
      expect(screen.getByText("Partly Cloudy")).toBeInTheDocument();
    });
    expect(screen.getByText("68°F")).toBeInTheDocument();
    expect(screen.getByText("moderate")).toBeInTheDocument();
  });

  it("calls fetch with correct coordinates", async () => {
    Object.defineProperty(navigator, "geolocation", {
      value: {
        getCurrentPosition: (success: PositionCallback) =>
          success({
            coords: { latitude: 39.7456, longitude: -97.0892 },
          } as GeolocationPosition),
      },
      writable: true,
    });

    (global.fetch as jest.Mock).mockResolvedValue({
      ok: true,
      json: () =>
        Promise.resolve({
          shortForecast: "Clear",
          temperature: 80,
          temperatureUnit: "F",
          temperatureCharacterization: "hot",
        }),
    });

    await act(async () => {
      render(<Map />);
    });

    expect(global.fetch).toHaveBeenCalledWith(
      "http://localhost:8080/api/forecast?lat=39.7456&lon=-97.0892"
    );
  });

  it("shows Unavailable when fetch fails", async () => {
    Object.defineProperty(navigator, "geolocation", {
      value: {
        getCurrentPosition: (success: PositionCallback) =>
          success({
            coords: { latitude: 40, longitude: -74 },
          } as GeolocationPosition),
      },
      writable: true,
    });

    (global.fetch as jest.Mock).mockRejectedValue(new Error("Network error"));

    await act(async () => {
      render(<Map />);
    });

    await waitFor(() => {
      expect(screen.getByText("Unavailable")).toBeInTheDocument();
    });
  });

  it("shows Unavailable when API returns non-ok response", async () => {
    Object.defineProperty(navigator, "geolocation", {
      value: {
        getCurrentPosition: (success: PositionCallback) =>
          success({
            coords: { latitude: 40, longitude: -74 },
          } as GeolocationPosition),
      },
      writable: true,
    });

    (global.fetch as jest.Mock).mockResolvedValue({
      ok: false,
      status: 502,
    });

    await act(async () => {
      render(<Map />);
    });

    await waitFor(() => {
      expect(screen.getByText("Unavailable")).toBeInTheDocument();
    });
  });

  it("displays correct weather icon for forecast", async () => {
    Object.defineProperty(navigator, "geolocation", {
      value: {
        getCurrentPosition: (success: PositionCallback) =>
          success({
            coords: { latitude: 40, longitude: -74 },
          } as GeolocationPosition),
      },
      writable: true,
    });

    (global.fetch as jest.Mock).mockResolvedValue({
      ok: true,
      json: () =>
        Promise.resolve({
          shortForecast: "Heavy Snow",
          temperature: 28,
          temperatureUnit: "F",
          temperatureCharacterization: "cold",
        }),
    });

    await act(async () => {
      render(<Map />);
    });

    await waitFor(() => {
      expect(screen.getByText("Heavy Snow")).toBeInTheDocument();
    });
    expect(screen.getByRole("img", { name: "Heavy Snow" })).toHaveTextContent(
      "❄️"
    );
  });
});
