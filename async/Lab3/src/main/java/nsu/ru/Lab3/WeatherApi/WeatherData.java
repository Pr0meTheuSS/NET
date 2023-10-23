package nsu.ru.Lab3.WeatherApi;
import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.fasterxml.jackson.annotation.JsonProperty;

@JsonIgnoreProperties(ignoreUnknown = true)
public class WeatherData {
    private Coord coord;
    private Weather[] weather;
    private String base;
    private Main main;
    private int visibility;
    private Wind wind;
    private Rain rain;
    private Clouds clouds;
    private long dt;
    private Sys sys;
    private int timezone;
    private int id;
    private String name;
    private int cod;

    public String getDescription() {
        return getWeather()[0].getDescription();
    }

    public double getTemp() {
        return getMain().getTemp();
    }

    public double getTempMin() {
        return getMain().getTempMin();
    }

    public double getTempMax() {
        return getMain().getTempMax();
    }

    public double getWindSpeed() {
        return getWind().getSpeed();
    }

    @JsonProperty("coord")
    public Coord getCoord() {
        return coord;
    }

    @JsonProperty("coord")
    public void setCoord(Coord coord) {
        this.coord = coord;
    }

    @JsonProperty("weather")
    public Weather[] getWeather() {
        return weather;
    }

    @JsonProperty("weather")
    public void setWeather(Weather[] weather) {
        this.weather = weather;
    }

    @JsonProperty("base")
    public String getBase() {
        return base;
    }

    @JsonProperty("base")
    public void setBase(String base) {
        this.base = base;
    }

    @JsonProperty("main")
    public Main getMain() {
        return main;
    }

    @JsonProperty("main")
    public void setMain(Main main) {
        this.main = main;
    }

    @JsonProperty("visibility")
    public int getVisibility() {
        return visibility;
    }

    @JsonProperty("visibility")
    public void setVisibility(int visibility) {
        this.visibility = visibility;
    }

    @JsonProperty("wind")
    public Wind getWind() {
        return wind;
    }

    @JsonProperty("wind")
    public void setWind(Wind wind) {
        this.wind = wind;
    }

    @JsonProperty("rain")
    public Rain getRain() {
        return rain;
    }

    @JsonProperty("rain")
    public void setRain(Rain rain) {
        this.rain = rain;
    }

    @JsonProperty("clouds")
    public Clouds getClouds() {
        return clouds;
    }

    @JsonProperty("clouds")
    public void setClouds(Clouds clouds) {
        this.clouds = clouds;
    }

    @JsonProperty("dt")
    public long getDt() {
        return dt;
    }

    @JsonProperty("dt")
    public void setDt(long dt) {
        this.dt = dt;
    }

    @JsonProperty("sys")
    public Sys getSys() {
        return sys;
    }

    @JsonProperty("sys")
    public void setSys(Sys sys) {
        this.sys = sys;
    }

    @JsonProperty("timezone")
    public int getTimezone() {
        return timezone;
    }

    @JsonProperty("timezone")
    public void setTimezone(int timezone) {
        this.timezone = timezone;
    }

    @JsonProperty("id")
    public int getId() {
        return id;
    }

    @JsonProperty("id")
    public void setId(int id) {
        this.id = id;
    }

    @JsonProperty("name")
    public String getName() {
        return name;
    }

    @JsonProperty("name")
    public void setName(String name) {
        this.name = name;
    }

    @JsonProperty("cod")
    public int getCod() {
        return cod;
    }

    @JsonProperty("cod")
    public void setCod(int cod) {
        this.cod = cod;
    }
}

class Coord {
    private double lon;
    private double lat;

    @JsonProperty("lon")
    public double getLon() {
        return lon;
    }

    @JsonProperty("lon")
    public void setLon(double lon) {
        this.lon = lon;
    }

    @JsonProperty("lat")
    public double getLat() {
        return lat;
    }

    @JsonProperty("lat")
    public void setLat(double lat) {
        this.lat = lat;
    }
}

class Weather {
    private int id;
    private String main;
    private String description;
    private String icon;

    // Геттеры и сеттеры

    @JsonProperty("id")
    public int getId() {
        return id;
    }

    @JsonProperty("id")
    public void setId(int id) {
        this.id = id;
    }

    @JsonProperty("main")
    public String getMain() {
        return main;
    }

    @JsonProperty("main")
    public void setMain(String main) {
        this.main = main;
    }

    @JsonProperty("description")
    public String getDescription() {
        return description;
    }

    @JsonProperty("description")
    public void setDescription(String description) {
        this.description = description;
    }

    @JsonProperty("icon")
    public String getIcon() {
        return icon;
    }

    @JsonProperty("icon")
    public void setIcon(String icon) {
        this.icon = icon;
    }
}

class Main {
    private double temp;
    private double feels_like;
    private double temp_min;
    private double temp_max;
    private int pressure;
    private int humidity;
    private int sea_level;
    private int grnd_level;

    // Геттеры и сеттеры

    @JsonProperty("temp")
    public double getTemp() {
        return temp;
    }

    @JsonProperty("temp")
    public void setTemp(double temp) {
        this.temp = temp;
    }

    @JsonProperty("feels_like")
    public double getFeelsLike() {
        return feels_like;
    }

    @JsonProperty("feels_like")
    public void setFeelsLike(double feels_like) {
        this.feels_like = feels_like;
    }

    @JsonProperty("temp_min")
    public double getTempMin() {
        return temp_min;
    }

    @JsonProperty("temp_min")
    public void setTempMin(double temp_min) {
        this.temp_min = temp_min;
    }

    @JsonProperty("temp_max")
    public double getTempMax() {
        return temp_max;
    }

    @JsonProperty("temp_max")
    public void setTempMax(double temp_max) {
        this.temp_max = temp_max;
    }

    @JsonProperty("pressure")
    public int getPressure() {
        return pressure;
    }

    @JsonProperty("pressure")
    public void setPressure(int pressure) {
        this.pressure = pressure;
    }

    @JsonProperty("humidity")
    public int getHumidity() {
        return humidity;
    }

    @JsonProperty("humidity")
    public void setHumidity(int humidity) {
        this.humidity = humidity;
    }

    @JsonProperty("sea_level")
    public int getSeaLevel() {
        return sea_level;
    }

    @JsonProperty("sea_level")
    public void setSeaLevel(int sea_level) {
        this.sea_level = sea_level;
    }

    @JsonProperty("grnd_level")
    public int getGrndLevel() {
        return grnd_level;
    }

    @JsonProperty("grnd_level")
    public void setGrndLevel(int grnd_level) {
        this.grnd_level = grnd_level;
    }
}

class Wind {
    private double speed;
    private int deg;
    private double gust;

    // Геттеры и сеттеры

    @JsonProperty("speed")
    public double getSpeed() {
        return speed;
    }

    @JsonProperty("speed")
    public void setSpeed(double speed) {
        this.speed = speed;
    }

    @JsonProperty("deg")
    public int getDeg() {
        return deg;
    }

    @JsonProperty("deg")
    public void setDeg(int deg) {
        this.deg = deg;
    }

    @JsonProperty("gust")
    public double getGust() {
        return gust;
    }

    @JsonProperty("gust")
    public void setGust(double gust) {
        this.gust = gust;
    }
}

class Rain {
    private double _1h;

    // Геттеры и сеттеры

    @JsonProperty("1h")
    public double get1h() {
        return _1h;
    }

    @JsonProperty("1h")
    public void set1h(double _1h) {
        this._1h = _1h;
    }
}

class Clouds {
    private int all;

    // Геттеры и сеттеры

    @JsonProperty("all")
    public int getAll() {
        return all;
    }

    @JsonProperty("all")
    public void setAll(int all) {
        this.all = all;
    }
}

class Sys {
    private int type;
    private int id;
    private String country;
    private long sunrise;
    private long sunset;

    // Геттеры и сеттеры

    @JsonProperty("type")
    public int getType() {
        return type;
    }

    @JsonProperty("type")
    public void setType(int type) {
        this.type = type;
    }

    @JsonProperty("id")
    public int getId() {
        return id;
    }

    @JsonProperty("id")
    public void setId(int id) {
        this.id = id;
    }

    @JsonProperty("country")
    public String getCountry() {
        return country;
    }

    @JsonProperty("country")
    public void setCountry(String country) {
        this.country = country;
    }

    @JsonProperty("sunrise")
    public long getSunrise() {
        return sunrise;
    }

    @JsonProperty("sunrise")
    public void setSunrise(long sunrise) {
        this.sunrise = sunrise;
    }

    @JsonProperty("sunset")
    public long getSunset() {
        return sunset;
    }

    @JsonProperty("sunset")
    public void setSunset(long sunset) {
        this.sunset = sunset;
    }
}
