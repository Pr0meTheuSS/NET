package nsu.ru.Lab3.Configs;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.boot.context.properties.ConfigurationProperties;
import org.springframework.stereotype.Component;

@Component
@ConfigurationProperties(prefix = "apis")
public class Configs {
    private String weatherApiKey;
    private String placesApiKey;
    private String locationsApiKey;

    public String getWeatherApiKey() {
        return weatherApiKey;
    }

    public void setWeatherApiKey(String weatherApiKey) {
        this.weatherApiKey = weatherApiKey;
    }

    public String getPlacesApiKey() {
        return placesApiKey;
    }

    public void setPlacesApiKey(String placesApiKey) {
        this.placesApiKey = placesApiKey;
    }

    public String getLocationsApiKey() {
        return locationsApiKey;
    }

    public void setLocationsApiKey(String locationsApiKey) {
        this.locationsApiKey = locationsApiKey;
    }
}
