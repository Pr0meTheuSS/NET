package nsu.ru.Lab3.controllers;

public class PlaceView {
    private int placeId;
    private String placeName;

    public PlaceView(int itemId, String placeName) {
        this.placeId = itemId;
        this.placeName = placeName;
    }

    public int getPlaceId() {
        return placeId;
    }

    public void setPlaceId(int itemId) {
        this.placeId = itemId;
    }

    public String getPlaceName() {
        return placeName;
    }

    public void setPlaceName(String placeName) {
        this.placeName = placeName;
    }
}
