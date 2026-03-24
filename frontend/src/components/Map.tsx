"use client";

import { useCallback, useEffect, useRef, useState } from "react";
import {
  MapContainer,
  TileLayer,
  Marker,
  Tooltip,
  useMap,
  useMapEvents,
} from "react-leaflet";
import L from "leaflet";
import "leaflet/dist/leaflet.css";

import markerIcon2x from "leaflet/dist/images/marker-icon-2x.png";
import markerIcon from "leaflet/dist/images/marker-icon.png";
import markerShadow from "leaflet/dist/images/marker-shadow.png";
import { getWeatherIcon } from "@/lib/weatherIcon";

L.Icon.Default.mergeOptions({
  iconUrl: markerIcon.src,
  iconRetinaUrl: markerIcon2x.src,
  shadowUrl: markerShadow.src,
});

const US_CENTER: L.LatLngExpression = [39.8283, -98.5795];
const DEFAULT_ZOOM = 4;
const LOCATION_ZOOM = 10;

interface Forecast {
  shortForecast: string;
  temperature: number;
  temperatureUnit: string;
  temperatureCharacterization: string;
}

function FlyToLocation({ position }: { position: L.LatLngExpression }) {
  const map = useMap();
  useEffect(() => {
    map.flyTo(position, LOCATION_ZOOM);
  }, [map, position]);
  return null;
}

function MapClickHandler({
  onMapClick,
}: {
  onMapClick: (lat: number, lon: number) => void;
}) {
  const draggedRef = useRef(false);

  useMapEvents({
    dragstart() {
      draggedRef.current = true;
    },
    click(e) {
      if (draggedRef.current) {
        draggedRef.current = false;
        return;
      }
      onMapClick(e.latlng.lat, e.latlng.lng);
    },
  });
  return null;
}

function ForecastMarker({
  position,
  forecast,
}: {
  position: [number, number];
  forecast: Forecast | null;
}) {
  const markerRef = useRef<L.Marker>(null);

  useEffect(() => {
    if (markerRef.current) {
      markerRef.current.openTooltip();
    }
  }, [position, forecast]);

  return (
    <Marker ref={markerRef} position={position}>
      <Tooltip
        direction="top"
        offset={[0, -30]}
        permanent={false}
        className="forecast-tooltip"
      >
        {!forecast ? (
          <span className="inline-flex items-center justify-center p-4">
            <svg
              className="h-6 w-6 animate-spin text-gray-500"
              xmlns="http://www.w3.org/2000/svg"
              fill="none"
              viewBox="0 0 24 24"
            >
              <circle
                className="opacity-25"
                cx="12"
                cy="12"
                r="10"
                stroke="currentColor"
                strokeWidth="4"
              />
              <path
                className="opacity-75"
                fill="currentColor"
                d="M4 12a8 8 0 018-8v4a4 4 0 00-4 4H4z"
              />
            </svg>
          </span>
        ) : (
          <div className="flex items-center gap-3 p-1" style={{ minWidth: 200 }}>
            <span className="text-4xl leading-none" role="img" aria-label={forecast.shortForecast}>
              {getWeatherIcon(forecast.shortForecast)}
            </span>
            <div className="flex flex-col">
              <span className="text-sm font-semibold">{forecast.shortForecast}</span>
              <span className="text-lg font-bold">
                {forecast.temperature}°{forecast.temperatureUnit}
              </span>
              <span className="text-xs text-gray-500 capitalize">
                {forecast.temperatureCharacterization}
              </span>
            </div>
          </div>
        )}
      </Tooltip>
    </Marker>
  );
}

export default function Map() {
  const [pinPosition, setPinPosition] = useState<[number, number] | null>(null);
  const [forecast, setForecast] = useState<Forecast | null>(null);
  const [initialFly, setInitialFly] = useState(false);

  const fetchForecast = useCallback((lat: number, lon: number) => {
    setForecast(null);
    const apiUrl = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";
    fetch(`${apiUrl}/api/forecast?lat=${lat}&lon=${lon}`)
      .then((res) => {
        if (!res.ok) throw new Error("Forecast unavailable");
        return res.json();
      })
      .then((data: Forecast) => setForecast(data))
      .catch(() =>
        setForecast({
          shortForecast: "Unavailable",
          temperature: 0,
          temperatureUnit: "F",
          temperatureCharacterization: "unknown",
        })
      );
  }, []);

  useEffect(() => {
    if (!navigator.geolocation) return;

    navigator.geolocation.getCurrentPosition(
      (pos) => {
        const lat = pos.coords.latitude;
        const lon = pos.coords.longitude;
        setPinPosition([lat, lon]);
        setInitialFly(true);
        fetchForecast(lat, lon);
      },
      () => {
        // Geolocation denied or unavailable — stay on default view
      }
    );
  }, [fetchForecast]);

  const handleMapClick = useCallback(
    (lat: number, lon: number) => {
      setPinPosition([lat, lon]);
      setInitialFly(false);
      fetchForecast(lat, lon);
    },
    [fetchForecast]
  );

  return (
    <MapContainer
      center={US_CENTER}
      zoom={DEFAULT_ZOOM}
      className="h-full w-full"
    >
      <TileLayer
        attribution='&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
        url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
      />
      <MapClickHandler onMapClick={handleMapClick} />
      {pinPosition && initialFly && <FlyToLocation position={pinPosition} />}
      {pinPosition && (
        <ForecastMarker position={pinPosition} forecast={forecast} />
      )}
    </MapContainer>
  );
}
