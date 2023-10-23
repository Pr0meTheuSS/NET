package nsu.ru.Lab3.WeatherApi;

import java.io.IOException;

public interface WeatherApiIface {
    WeatherData getWeatherAtPoint(String lat, String lng)  throws IOException, InterruptedException;
}
