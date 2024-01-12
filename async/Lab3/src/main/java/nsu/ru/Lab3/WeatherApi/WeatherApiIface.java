package nsu.ru.Lab3.WeatherApi;

import java.util.concurrent.CompletableFuture;

import java.io.IOException;

public interface WeatherApiIface {
    CompletableFuture<WeatherData> getWeatherAtPoint(String lat, String lng)  throws IOException, InterruptedException;
}
