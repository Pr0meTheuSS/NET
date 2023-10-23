package nsu.ru.Lab3.LocationApi;
import com.fasterxml.jackson.annotation.JsonProperty;
import java.util.List;

public class LocationResponseDTO {


    public String getHitsLat(int id) {
        return getHits().get(id).getPoint().getLat();
    }

    public String getHitsLon(int id) {
        return getHits().get(id).getPoint().getLng();
    }

    @JsonProperty("hits")
    private Location[] hits;

    @JsonProperty("locale")
    private String locale;

    public List <Location> getHits() {
        return List.of(hits);
    }

    public void setHits(Location[] hits) {
        this.hits = hits;
    }

    public String getLocale() {
        return locale;
    }

    public void setLocale(String locale) {
        this.locale = locale;
    }

}
