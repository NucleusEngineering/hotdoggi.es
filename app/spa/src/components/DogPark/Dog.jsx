export default function Dog(props) {
  const { type, location, name } = { ...props };

  return (
    <div
      style={{
        display: "flex",
        flexDirection: "column",
        alignItems: "center",
        gridColumn: location.x + 1, // css grids start from index 1 :(
        gridRow: location.y + 1,
      }}
    >
      <img height="75" width="75" src={`dogs/${type}.png`} />
      <span>{name}</span>
    </div>
  );
}
