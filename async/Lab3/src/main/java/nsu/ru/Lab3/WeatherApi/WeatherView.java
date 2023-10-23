package nsu.ru.Lab3.WeatherApi;


public class WeatherView {
    private String description;
    private String temperature;
    private String minTemperature;
    private String maxTemperature;
    private String windSpeed;

    public String getDescription() {
        return description;
    }
    public void setDescription(String description) {
        this.description = description;
    }
    public String getTemperature() {
        return temperature;
    }
    public void setTemperature(String temperature) {
        this.temperature = temperature;
    }
    public String getMinTemperature() {
        return minTemperature;
    }
    public void setMinTemperature(String minTemperature) {
        this.minTemperature = minTemperature;
    }
    public String getMaxTemperature() {
        return maxTemperature;
    }
    public void setMaxTemperature(String maxTemperature) {
        this.maxTemperature = maxTemperature;
    }
    public String getWindSpeed() {
        return windSpeed;
    }
    public void setWindSpeed(String windSpeed) {
        this.windSpeed = windSpeed;
    }
}