package nsu.ru.Lab3.PlacesApi;

public class PlaceInfoView {
    private String name;
    private Address address;
    private String kinds;

    public String getName() {
        return name;
    }

    public void setName(String name) {
        this.name = name;
    }

    public String getAddress() {
        if (address == null) {
            return "";
        }

        String ret = address.getCity() + " " + address.getState();
        if (address.getRoad() != null && address.getRoad() != "") {
            ret +=  " " + address.getRoad() + " ";
        }
        if (address.getHouse_number() != null && address.getHouse_number() != "") {
            ret += " " + address.getHouse_number() + " ";
        }

        return ret.replace("  ", " ");
    }

    public void setAddress(Address address) {
        this.address = address;
    }

    public String getKinds() {
        return kinds;
    }

    public void setKinds(String kinds) {
        this.kinds = kinds;
    }
}
